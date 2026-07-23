package limiterCache

import (
	"context"
	_ "embed"
	"fmt"
	"gateway-api/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed luaScripts/bucketLimiter.lua
var tokenBucketLua string

type Cache struct {
	client            *redis.Client
	tokenBucketScript *redis.Script
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

	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{
		client:            client,
		tokenBucketScript: redis.NewScript(tokenBucketLua),
	}, nil
}

func (c *Cache) Take(
	ctx context.Context,
	key string,
	capacity int64,
	refillRate float64,
	requested int64,
	ttl time.Duration,
) (bool, int64, error) {
	const op = "cache.redis.NewClient"

	result, err := c.tokenBucketScript.Run(
		ctx,
		c.client,
		[]string{key},
		capacity,
		refillRate,
		time.Now().Unix(),
		requested,
		int64(ttl.Seconds()),
	).Result()

	if err != nil {
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	values := result.([]interface{})

	return values[0].(int64) == 1, values[1].(int64), nil
}
