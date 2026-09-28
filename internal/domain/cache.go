package domain

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (any, error)
	Del(ctx context.Context, key string) error
}
