package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

func (m *Middlewares) LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := WrapResponseWriter(w)
		reqID := ulid.Make().String()

		m.logger.Info("[START]",
			slog.String("ReqID", reqID),
			slog.String("Method", r.Method),
			slog.String("URL", r.URL.String()),
			slog.String("RemoteAddr", r.RemoteAddr),
			slog.String("UserAgent", r.UserAgent()),
			// slog.Any("Headers", r.Header),
		)

		defer func() {
			result := "failed"
			duration := time.Since(rw.Start)

			if rw.Finished {
				result = "successfully"
			}

			m.logger.Info("[END]",
				slog.String("ReqID", reqID),
				slog.Int("Status", rw.StatusCode),
				slog.Duration("Duration", duration),
				slog.String("Result", result),
			)
		}()

		next.ServeHTTP(rw, r)
	}
}
