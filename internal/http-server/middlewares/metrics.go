package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/metrics"
)

func routeLabel(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}

	segments := strings.Split(strings.Trim(trimmed, "/"), "/")
	for i, segment := range segments {
		if _, err := strconv.Atoi(segment); err == nil {
			segments[i] = ":id"
		}
	}

	return "/" + strings.Join(segments, "/")
}

func (m *Middlewares) MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := WrapResponseWriter(w)
		route := routeLabel(r.URL.Path)

		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		next.ServeHTTP(rw, r)

		statusCode := rw.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		metrics.RequestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(statusCode)).Inc()

		switch statusCode / 100 {
		case 4:
			metrics.ClientErrorsTotal.WithLabelValues(r.Method, route).Inc()
		case 5:
			metrics.ServerErrorsTotal.WithLabelValues(r.Method, route).Inc()
		default:
		}

		metrics.RequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(rw.Start).Seconds())
	}
}
