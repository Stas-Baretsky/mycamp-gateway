package goredis

import (
	"context"
	"fmt"
	"gateway-api/internal/config"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewClient(ctx context.Context, cfg config.Config) (*Cache, error) {
	const op = "cache.redis.NewClient"

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Cache.Timeout,
		WriteTimeout: cfg.Cache.Timeout,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{
		client: client,
	}, nil
}

func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) {

}
