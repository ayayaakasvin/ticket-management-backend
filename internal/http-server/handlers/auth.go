package handlers

import (
	"net/http"
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
			h.logger.Errorf("failed to set sessionId: %v", sessionId)
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return 
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}


func (h *Handlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerReq request.UserRequest
		if err := bindjson.BindJson(r.Body, &registerReq); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
			return
		}

		if !(validinput.IsValidPassword(registerReq.Password) && validinput.IsValidUsername(registerReq.Username)) {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid credentials for register")
			return
		}

		hashed, err := bcrypthashing.BcryptHashing(registerReq.Password)
		if err != nil {
			h.logger.Errorf("Bcrypt error during hashing: %v", err)
			response.SendErrorJson(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if err := h.userRepo.RegisterUser(r.Context(), registerReq.Username, hashed); err != nil {
			h.logger.Errorf("RegisterUser error: %v", err)
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to register")
			return
		}

		response.SendSuccessJson(w, http.StatusCreated, nil)
	}
}

func (h *Handlers) LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id, ok := r.Context().Value(ctx.CtxSessionIDKey).(string)
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "missing session id")
			return
		}

		if err := h.cache.Del(r.Context(), session_id); err != nil {
			h.logger.Errorf("failed to delete session id: %v", err)
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
			h.logger.Errorf("failed to set session id: %v", err)
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}
