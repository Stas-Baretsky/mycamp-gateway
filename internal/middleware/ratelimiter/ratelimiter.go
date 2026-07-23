package ratelimiter

import (
	"context"
	"errors"
	"gateway-api/internal/config"
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
	Take(
		ctx context.Context,
		key string,
		capacity int64,
		refillPerSecond float64,
		requested int64,
		ttl time.Duration,
	) (bool, int64, error)
}

// /TASK :: протянуть контекст и сделать greasful shutdown
func NewRateLimiter(
	ctx context.Context,
	log *slog.Logger,
	storage RateLimiterStorage,
	cfg *config.RateLimiter,
) (*RateLimiter, error) {
	log.Info("Init new ratelimiter")
	return &RateLimiter{
		log:     log,
		storage: storage,
	}, nil
}

func (rl *RateLimiter) Approve() (bool, error) {
	return true, nil
}

func (rl *RateLimiter) Stop(ctx context.Context) {

}
