package redis

import (
	"context"
	"fmt"
	"gateway-api/internal/cache/redis"
	"gateway-api/internal/config"

	"github.com/redis/go-redis"
)

type Cache struct {
	client *redis.Client
}

func NewClient(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	const op = "cache.redis.NewClient"

	db, err := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &db, nil
}
