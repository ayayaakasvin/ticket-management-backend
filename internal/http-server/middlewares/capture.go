package middlewares

import "net/http"

func (m *Middlewares) CaptureMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := WrapResponseWriter(w)

		next.ServeHTTP(rw, r)
	}
}
