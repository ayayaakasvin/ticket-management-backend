package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/token"

	"github.com/google/uuid"
)

const (
	AuthorizationHeader  = "Authorization"
	sessionKeyPrefix     = "session:"
	userSessionKeyPrefix = "user_session:"
)

func sessionKey(sessionID string) string {
	return sessionKeyPrefix + sessionID
}

func userSessionKey(userID int64) string {
	return userSessionKeyPrefix + strconv.FormatInt(userID, 10)
}

func (h *Handlers) storeSession(ctx context.Context, userID int64, sessionID string) error {
	if oldSession, err := h.cc.Get(ctx, userSessionKey(userID)); err == nil {
		if oldSessionID, ok := oldSession.(string); ok && oldSessionID != "" {
			_ = h.cc.Del(ctx, sessionKey(oldSessionID))
		}
	}

	if err := h.cc.Set(ctx, sessionKey(sessionID), true, token.AccessTokenTTL); err != nil {
		return err
	}

	if err := h.cc.Set(ctx, userSessionKey(userID), sessionID, token.RefreshTokenTTL); err != nil {
		_ = h.cc.Del(ctx, sessionKey(sessionID))
		return err
	}

	return nil
}

// Login authenticates a user and returns an access token and refresh-token cookie.
// @Summary      Log in
// @Description  Authenticates a user and returns an access token. The refresh token is set as an HttpOnly cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      domain.UserRequest  true  "Login payload"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/auth/login [post]
func (h *Handlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginReq domain.UserRequest
		if err := helper.BindJson(r.Body, &loginReq); err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "request bind error"})
			return
		}

		userObj, err := h.ur.GetByUsername(r.Context(), loginReq.Username)
		if err != nil {
			h.logger.Error("GetByUsername user repository error", "err", err, "username", loginReq.Username)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		} else if userObj == nil {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		if err := helper.ComparePasswordAndHash(loginReq.Password, userObj.PasswordHash); err != nil {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		sessionId := helper.MustUUIDV7().String()

		accessToken, err := h.jwtM.GenerateToken(token.NewAccessTokenClaims(userObj.ID, sessionId, token.AccessTokenTTL))
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "access token generation error"})
			h.logger.Error("token generation error", "err", err)
			return
		}

		refreshToken, err := h.jwtM.GenerateToken(token.NewRefreshTokenClaims(userObj.ID, token.RefreshTokenTTL))
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "refresh token generation error"})
			h.logger.Error("token generation error", "err", err)
			return
		}

		if err := h.storeSession(r.Context(), userObj.ID, sessionId); err != nil {
			h.logger.Error("failed to store session", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "cache error"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/api/auth",
			HttpOnly: true,                 // for local testing, set to true in production
			Secure:   false,                // for local testing, set to true in production
			SameSite: http.SameSiteLaxMode, // for local testing, set to SameSiteStrictMode in production
			MaxAge:   int(token.RefreshTokenTTL.Seconds()),
		})

		data := map[string]any{}
		data["access-token"] = accessToken

		if err := h.cc.Set(r.Context(), sessionId, true, token.AccessTokenTTL); err != nil {
			h.logger.Error("failed to set session id", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "cache error"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// Register creates a user account.
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      domain.UserRequest  true  "Register payload"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/auth/register [post]
func (h *Handlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerReq domain.UserRequest
		if err := helper.BindJson(r.Body, &registerReq); err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "request bind error"})
			return
		}

		if err := helper.IsValidUsername(registerReq.Username); err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "username must be at least 3 characters and contain a letter"})
			return
		}
		if err := helper.IsValidPassword(registerReq.Password); err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "password must be at least 8 characters and contain uppercase, lowercase, and numeric characters"})
			return
		}

		hashed, err := helper.BcryptHashing(registerReq.Password)
		if err != nil {
			h.logger.Error("bcrypt hashing failed", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		}

		userObj, err := h.ur.GetByUsername(r.Context(), registerReq.Username)
		if err != nil {
			h.logger.Error("GetByUsername user repository error", "err", err, "username", userObj.Username)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		} else if userObj != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "username already exists"})
			return
		}

		id, err := h.ur.Create(r.Context(), &domain.User{
			Username:     registerReq.Username,
			PasswordHash: hashed,
			Role:         domain.Client,
		})

		if err != nil {
			h.logger.Error("register user failed", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to register"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusCreated, map[string]any{
			"message": "user registered successfully",
			"user_id": id,
		})
	}
}

// Logout invalidates the authenticated session and clears the refresh-token cookie.
// @Summary      Log out
// @Description  Invalidates the current session (deletes session id)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/auth/logout [post]
func (h *Handlers) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id, ok := r.Context().Value(ctx.CtxSessionIDKey).(string)
		if !ok {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "session id not found in context"})
			return
		}

		userID, ok := r.Context().Value(ctx.CtxUserIDKey).(int64)
		if ok && userID != 0 {
			_ = h.cc.Del(r.Context(), userSessionKey(userID))
		}

		if err := h.cc.Del(r.Context(), session_id); err != nil {
			h.logger.Error("failed to delete session id", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete session id"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/api/auth",
			HttpOnly: true,                 // for local testing, set to true in production
			Secure:   false,                // for local testing, set to true in production
			SameSite: http.SameSiteLaxMode, // for local testing, set to SameSiteStrictMode in production
			MaxAge:   -1,                   // Delete the cookie
			Expires:  time.Unix(0, 0),
		})

		helper.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
	}
}

// Refresh exchanges the refresh-token cookie for a new access token.
// @Summary      Refresh access token
// @Description  Exchanges the refresh_token cookie for a new access token. Send the refresh_token as a cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200            {object}  map[string]interface{}
// @Failure      401            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /api/auth/refresh [post]
func (h *Handlers) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "refresh token cookie not found"})
			return
		}

		refreshTokenString := cookie.Value
		cl, err := h.jwtM.Validate(refreshTokenString, &token.RefreshTokenClaims{})
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "failed to validate jwt"})
			return
		}

		fullClaims, ok := cl.(*token.RefreshTokenClaims)
		if !ok {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
			return
		}

		if fullClaims.UserID == 0 {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "user_id is missing in refresh token"})
			return
		}

		sessionId := uuid.New().String()

		accessToken, err := h.jwtM.GenerateToken(token.NewAccessTokenClaims(fullClaims.UserID, sessionId, token.AccessTokenTTL))
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
			h.logger.Error("failed to generate token")
			return
		}

		data := map[string]any{}
		data["access-token"] = accessToken

		if err := h.storeSession(r.Context(), fullClaims.UserID, sessionId); err != nil {
			h.logger.Error("failed to store session", "err", err)
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to set session id"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}
