package valkey

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/redis/go-redis/v9"
)

const origin = "Redis/Valkey"

// for storing methods of storing and retrieving session_id
type Valkey struct {
	connection *redis.Client
}

func NewValkeyWithURL(url string) (*Valkey, error) {
	ctx := context.Background()
	opt, err := redis.ParseURL(url)
	log.Printf("URL: %s", url)
	if err != nil {
		return nil,  fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	// for latency
	opt.DialTimeout = 30 * time.Second // Increased for Singapore region
	opt.ReadTimeout = 30 * time.Second
	opt.WriteTimeout = 30 * time.Second
	opt.PoolSize = 10
	opt.PoolTimeout = 30 * time.Second

	conn := redis.NewClient(opt)
	if err := conn.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v\n", err)
	}

	return &Valkey{
		connection: conn,
	}, nil
}

func (c *Valkey) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.connection.Set(ctx, key, value, ttl).Err()
}

func (c *Valkey) Get(ctx context.Context, key string) (any, error) {
	return c.connection.Get(ctx, key).Result()
}

func (c *Valkey) Del(ctx context.Context, key string) error {
	return c.connection.Del(ctx, key).Err()
}

func (c *Valkey) SetNX(ctx context.Context, key string, value any, ttl time.Duration) *redis.BoolCmd {
	return c.connection.SetNX(ctx, key, value, ttl)
}
