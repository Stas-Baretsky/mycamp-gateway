package ratelimiter

import (
	"context"
	"errors"
	"gateway-api/internal/config"
	"gateway-api/internal/utils/bucketlimiter"
	"log/slog"
	"time"
)

var (
	KeyDoesNotExist = errors.New("key does not exist")
)

type RateLimiter struct {
	log     *slog.Logger
	storage RateLimiterStorage
}

type RateLimiterStorage interface {
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration)
	Read(ctx context.Context, key string) (int64, error)
}

// /TASK :: протянуть контекст и сделать greasful shutdown
func NewRateLimiter(
	ctx context.Context,
	log *slog.Logger,
	storage *RateLimiterStorage,
	cfg *config.RateLimiter,
) (*RateLimiter, error) {

	rLimiter := bucketlimiter.MakeBucketLimiter()
	rLimiter.Init(
		storage,
		&bucketlimiter.LimiterSettings{
			Speed:      cfg.FillingSpeed,
			Limit:      cfg.BucketCap,
			CallWeight: cfg.ReqWeight,
			TokenLife:  cfg.TokenLife,
		})

	return &RateLimiter{}, nil
}
