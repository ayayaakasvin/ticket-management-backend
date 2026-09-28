package middlewares

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
)

func (m *Middlewares) AdminOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := ctx.GetUserIDFromContext(r.Context())
		if !ok {
			m.logger.Warn("missing user id in request context")
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		ok = m.ur.IsAdmin(r.Context(), userID)
		if !ok {
			m.logger.Warn("user tried to access protected admin route", "user_id", userID)
			helper.WriteJSONResponse(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			return
		}

		next.ServeHTTP(w, r)
	}
}