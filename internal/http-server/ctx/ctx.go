package ctx

import (
	"context"
	"net/http"
	"strconv"
)

// used as key for context.WithValue
type contextKey string

// most used keys in ctx
const (
	CtxSessionIDKey contextKey = "session_id"
	CtxUserIDKey    contextKey = "user_id"
)

func WrapValueIntoRequest(r *http.Request, key contextKey, value any) *http.Request {
	newCtx := context.WithValue(r.Context(), key, value)
	return r.WithContext(newCtx)
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	switch v := ctx.Value(CtxUserIDKey).(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint64:
		return int64(v), true
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return parsed, true
		}
	}

	return 0, false
}
