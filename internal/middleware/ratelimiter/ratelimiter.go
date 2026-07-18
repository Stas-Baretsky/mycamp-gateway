package ratelimiter

import (
	"context"
	"gateway-api/internal/config"
	"gateway-api/internal/utils/limiter"
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

func NewRateLimiter(
	ctx context.Context,
	log *slog.Logger,
	storage RateLimiterStorage,
	cfg config.Config,
) (RateLimiter, error) {

	limiter := limiter.MakeBucketLimiter()
	limiter.Init(
		cfg.RaceLimiter.FillingSpeed,
		cfg.RaceLimiter.BucketCap,
		cfg.RaceLimiter.ReqWeight,
		cfg.RaceLimiter.TokenLife,
	)

	return RateLimiter{}, nil
}
