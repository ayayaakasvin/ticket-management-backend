package middlewares

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
)


func (h *Middlewares) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		if !h.rlm.GetLimiter(userID).Allow() {
			http.Error(w, "Too many uploads, slow down", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}