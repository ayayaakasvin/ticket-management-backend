package middlewares

import (
	"net/http"
)

// CORS middleware
func (mw *Middlewares) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mw.writeHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (mw *Middlewares) writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", mw.allowedOrigins)
	w.Header().Set("Access-Control-Allow-Methods", mw.allowedMethods)
	w.Header().Set("Access-Control-Allow-Headers", mw.allowedHeaders)

	if mw.allowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}
