package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/bcrypthashing"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/bindjson"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/jwttool"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/validinput"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/request"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/response"
	"github.com/ayayaakasvin/oneflick-ticket/internal/repo/postgresql"

	"github.com/google/uuid"
)

var expTimeAccessToken time.Duration = time.Minute * 15
var expTimeRefreshToken time.Duration = time.Hour * 168

const (
	AuthorizationHeader = "Authorization"
	RegisterForm        = "%s:%s:%s" //
	RegisterTTL         = time.Minute * 10
)

func (h *Handlers) LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginReq request.UserRequest
		if err := bindjson.BindJson(r.Body, &loginReq); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
			return
		}

		userId, err := h.userRepo.AuthentificateUser(r.Context(), loginReq.Username, loginReq.Password)
		if err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				response.SendErrorJson(w, http.StatusUnauthorized, "invalid credentials")
			case postgresql.UnAuthorized:
				response.SendErrorJson(w, http.StatusUnauthorized, "invalid credentials")
			}
			return
		}

		sessionId := uuid.New().String()
		accessToken := jwttool.GenerateAccessToken(userId, sessionId, expTimeAccessToken)
		refreshToken := jwttool.GenerateRefreshToken(userId, expTimeRefreshToken)

		data := response.NewData()
		data["access-token"] = accessToken
		data["refresh-token"] = refreshToken
		h.logger.Info(data)

		if err := h.cache.Set(r.Context(), sessionId, true, expTimeAccessToken); err != nil {
			h.logger.WithField("session_id", sessionId).Error("failed to set session id")
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

// func (h *Handlers) Register() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var registerReq request.UserRequest
// 		if err := bindjson.BindJson(r.Body, &registerReq); err != nil {
// 			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
// 			return
// 		}

// 		if !(validinput.IsValidPassword(registerReq.Password) && validinput.IsValidUsername(registerReq.Username)) {
// 			response.SendErrorJson(w, http.StatusBadRequest, "invalid credentials for register")
// 			return
// 		}

// 		hashed, err := bcrypthashing.BcryptHashing(registerReq.Password)
// 		if err != nil {
// 			h.logger.WithError(err).Error("bcrypt hashing failed")
// 			response.SendErrorJson(w, http.StatusInternalServerError, "Internal Server Error")
// 			return
// 		}

// 		if err := h.userRepo.RegisterUser(r.Context(), registerReq.Username, hashed, registerReq.Email); err != nil {
// 			h.logger.WithError(err).Error("register user failed")
// 			response.SendErrorJson(w, http.StatusInternalServerError, "failed to register")
// 			return
// 		}

// 		response.SendSuccessJson(w, http.StatusCreated, nil)
// 	}
// }

func (h *Handlers) RegisterStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerReq request.UserRequest
		if err := bindjson.BindJson(r.Body, &registerReq); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
			return
		}

		if !validinput.ValidateUserRequest(&registerReq) {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid credentials for register")
			return
		}

		code := h.smtp.GenerateRandomSequence()
		codeString := strconv.Itoa(code)

		registerForm := convertFormToString(&registerReq)

		if err := h.cache.Set(r.Context(), codeString, registerForm, RegisterTTL); err != nil {
			h.logger.WithError(err).Error("failed to set user request in cache")
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		if err := h.smtp.SendCode("Verification code", code, []string{registerReq.Email}); err != nil {
			h.logger.WithError(err).WithField("email", registerReq.Email).Error("failed to send verification code to email")
			response.SendErrorJson(w, http.StatusInternalServerError, "smtp error")
			return
		}

		response.SendSuccessJson(w, http.StatusAccepted, nil)
	}
}

func (h *Handlers) RegisterVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.RegisterVerify
		if err := bindjson.BindJson(r.Body, &req); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid request format")
			return
		}

		// Get cached form
		formAny, err := h.cache.Get(r.Context(), strconv.Itoa(req.Code))
		if err != nil {
			h.logger.WithError(err).Warn("failed to get cached form")
			response.SendErrorJson(w, http.StatusBadRequest, "verification failed")
			return
		}

		formString, ok := formAny.(string)
		if !ok || formString == "" {
			h.logger.Warn("cached form is empty or invalid")
			response.SendErrorJson(w, http.StatusBadRequest, "verification failed")
			return
		}

		form, err := fetchFormFromString(formString)
		if err != nil {
			h.logger.WithError(err).Warn("failed to parse cached form")
			response.SendErrorJson(w, http.StatusBadRequest, "verification failed")
			return
		}

		hashed, err := bcrypthashing.BcryptHashing(form.Password)
		if err != nil {
			h.logger.WithError(err).Error("failed to hash password")
			response.SendErrorJson(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if err := h.userRepo.RegisterUser(r.Context(), form.Username, hashed, form.Email); err != nil {
			h.logger.WithError(err).Warn("failed to register user")
			// Do NOT say "username already exists"
			response.SendErrorJson(w, http.StatusBadRequest, "registration failed")
			return
		}

		response.SendSuccessJson(w, http.StatusCreated, nil)
	}
}

// TODO: Register start and verify to implement

func (h *Handlers) LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id, ok := r.Context().Value(ctx.CtxSessionIDKey).(string)
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "missing session id")
			return
		}

		if err := h.cache.Del(r.Context(), session_id); err != nil {
			h.logger.WithError(err).Error("failed to delete session id")
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

func (h *Handlers) RefreshTheToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			response.SendErrorJson(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		refreshTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if refreshTokenString == authHeader {
			response.SendErrorJson(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		claims, err := jwttool.ValidateJWT(refreshTokenString)
		if err != nil {
			response.SendErrorJson(w, http.StatusUnauthorized, "failed to validate jwt")
			return
		}

		userIdAny, ok := claims["user_id"]
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "user_id is missing in refresh token")
			return
		}

		userId, err := jwttool.FetchUserID(userIdAny)
		if err != nil {
			response.SendErrorJson(w, http.StatusUnauthorized, "user_id is invalid")
			return
		}

		sessionId := uuid.New().String()
		accessToken := jwttool.GenerateAccessToken(userId, sessionId, expTimeAccessToken)

		data := response.NewData()
		data["access-token"] = accessToken
		h.logger.Info(data)

		if err := h.cache.Set(r.Context(), sessionId, true, expTimeAccessToken); err != nil {
			h.logger.WithError(err).Error("failed to set session id in cache")
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

func fetchFormFromString(form string) (*request.UserRequest, error) {
	formSlice := strings.Split(form, ":")

	if len(formSlice) == 0 {
		return nil, fmt.Errorf("empty form string")
	}

	var req *request.UserRequest
	req.Username = formSlice[0]
	req.Password = formSlice[1]
	req.Email = formSlice[2]

	return req, nil
}

func convertFormToString(req *request.UserRequest) string {
	return fmt.Sprintf(RegisterForm, req.Username, req.Password, req.Email)
}
