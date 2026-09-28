package middlewares

import (
	"net/http"
	"runtime/debug"
)

// recover middleware
func (mw *Middlewares) RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				mw.logger.Error(
					"panic recovered",
					"panic", recovered,
					"stack", string(debug.Stack()),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}
