package gcache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

type GCache struct {
	cc *cache.Cache
}

// Close implements [domain.Cache].
func (g *GCache) Close() error {
	return nil
}

// Del implements [domain.Cache].
func (g *GCache) Del(ctx context.Context, key string) error {
	g.cc.Delete(key)
	return nil
}

// Get implements [domain.Cache].
func (g *GCache) Get(ctx context.Context, key string) (any, error) {
	val, found := g.cc.Get(key)
	if !found {
		return nil, errors.New("key not found")
	}
	return val, nil
}

// Set implements [domain.Cache].
func (g *GCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	g.cc.Set(key, value, ttl)
	return nil
}

// SetNX implements [domain.Cache].
func (g *GCache) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	g.cc.Set(key, value, ttl)
	return true, nil
}

func NewGCache(access time.Duration) *GCache {
	return &GCache{
		cache.New(access, cache.NoExpiration),
	}
}