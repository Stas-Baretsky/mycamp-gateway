package ratelimiter

import (
	"context"
	"errors"
	"gateway-api/internal/config"
	"gateway-api/internal/utils/bucketlimiter"
	"log/slog"
)

var (
	KeyDoesNotExist = errors.New("key does not exist")
)

type RateLimiter struct {
	log     *slog.Logger
	storage RateLimiterStorage
}

type RateLimiterStorage interface {
	Take(
		ctx context.Context,
		key string,
		capacity int64,
		refillRate int64,
		///Deadline
	) (bool, int64, error)
}

// /TASK :: протянуть контекст и сделать greasful shutdown
func NewRateLimiter(
	ctx context.Context,
	log *slog.Logger,
	storage RateLimiterStorage,
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
