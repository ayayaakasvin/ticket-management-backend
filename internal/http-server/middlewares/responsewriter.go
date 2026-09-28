package middlewares

import (
	"net/http"
	"time"
)

// rw implementation for tracking if request was handled successfully, using bool value rw.finished and assigning true if WriteHeader was called.
// Looks like this Req -> Logger -> Handler -> Logger (checks rw.finished value) -> based on it Result shows up.
type ResponseWriter struct {
	http.ResponseWriter

	Start      time.Time
	StatusCode int
	Finished   bool
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if rw.Finished {
		return
	}

	rw.StatusCode = code
	rw.Finished = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.Finished {
		rw.StatusCode = http.StatusOK
		rw.Finished = true
		rw.ResponseWriter.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,

		Start:    time.Now(),
		Finished: false,
	}
}

func WrapResponseWriter(w http.ResponseWriter) *ResponseWriter {
	if rw, ok := w.(*ResponseWriter); ok {
		return rw
	}

	return NewResponseWriter(w)
}
