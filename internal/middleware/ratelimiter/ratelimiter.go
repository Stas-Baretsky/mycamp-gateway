package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"gateway-api/internal/config"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	KeyDoesNotExist = errors.New("key does not exist")
)

type RateLimiter struct {
	log     *slog.Logger
	storage RateLimiterStorage
	cfg     config.RateLimiter
}

type RateLimiterStorage interface {
	Take(
		ctx context.Context,
		key string,
		capacity int64,
		refillPerSecond float64,
		requested int64,
		ttl time.Duration,
	) (bool, error)
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
		cfg:     *cfg,
	}, nil
}

func (rl *RateLimiter) Approve(
	ctx context.Context,
	key string,
) (bool, error) {
	const op = "middleware.ratelimiter.ratelimiter.Approve"
	approve, err := rl.storage.Take(
		ctx,
		key,
	)
	if err != nil {
		log.Error("Redis err", sl.Err(fmt.Errorf("%s: %w", op, err)))
	}
	return approve, nil
}

func (rl *RateLimiter) Stop(ctx context.Context) {

}
