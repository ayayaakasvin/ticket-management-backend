package valkey

import (
	"context"
	"fmt"
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
	conn := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Redis server address
			Password: cfg.Password,                             // No password set
			Username: cfg.User,
			DB:       0, // Default DB)
		},
	)

	if err := conn.Ping(context.Background()).Err(); err != nil {
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
