package ratelimiter

import (
	"context"
	"log/slog"
	"time"
)

type RateLimiter struct {
	log     *slog.Logger
	storage RateLimiterStorage
}

type RateLimiterStorage interface {
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration)
}

func NewRateLimiter() (RateLimiter, error) {
	///Написать rate limiter и пометить в utils
	return RateLimiter{}, nil
}
