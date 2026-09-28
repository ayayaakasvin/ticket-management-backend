package middlewares

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerMiddlewareHandlesStandardResponseWriter(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelInfo}))

	mw := &Middlewares{logger: logger}

	handler := mw.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	logOutput := logs.String()
	if !strings.Contains(logOutput, "[END]") {
		t.Fatalf("expected logger output to contain request completion, got: %s", logOutput)
	}
}

func TestRecoverMiddlewareCallsNextHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mw := &Middlewares{logger: logger}

	handler := mw.RecoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
}

func TestRecoverMiddlewareReturns500OnPanic(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mw := &Middlewares{logger: logger}

	handler := mw.RecoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
