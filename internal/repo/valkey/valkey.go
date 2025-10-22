package valkey

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/redis/go-redis/v9"
)

const origin = "Redis/Valkey"

// for storing methods of storing and retrieving session_id
type Cache struct {
	connection *redis.Client
}

func NewValkeyClient(cfg config.RedisConfig, shutdownChan inner.ShutdownChannel) inner.Cache {
	ctx := context.Background()
	opt, err := redis.ParseURL(cfg.URL)
	log.Printf("URL: %s", cfg.URL)
	if err != nil {
		msg := fmt.Sprintf("failed to parse Redis URL: %v", err)
		shutdownChan.Send(inner.ShutdownMessage, origin, msg)
		return nil
	}

	// for latency
	opt.DialTimeout = 30 * time.Second // Increased for Singapore region
    opt.ReadTimeout = 30 * time.Second
    opt.WriteTimeout = 30 * time.Second
    opt.PoolSize = 10
    opt.PoolTimeout = 30 * time.Second

	conn := redis.NewClient(opt)
	if err := conn.Ping(ctx).Err(); err != nil {
		msg := fmt.Sprintf("failed to connect to db: %v\n", err)
		shutdownChan.Send(inner.ShutdownMessage, origin, msg)
		return nil
	}

	return &Cache{
		connection: conn,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.connection.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	return c.connection.Get(ctx, key).Result()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.connection.Del(ctx, key).Err()
}

func (c *Cache) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	return c.connection.SetNX(ctx, key, value, ttl)
}
