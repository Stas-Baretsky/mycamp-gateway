package bootstrap

import (
	"context"
	"fmt"

	"gateway-api/internal/app"
	limiterCache "gateway-api/internal/cache/limiterCache"
	"gateway-api/internal/config"
	"gateway-api/internal/middleware/ratelimiter"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Build(ctx context.Context, cfg config.Config) (*app.Application, error) {
	const op = "bootstrap.bootstrap.Build"

	log := setupLogger(cfg.Env)
	log.Info("Starting service...", slog.String("env", cfg.Env))

	redisStorage, err := limiterCache.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	limiter, err := ratelimiter.NewRateLimiter(ctx, log, redisStorage, &cfg.RateLimiter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &app.Application{
		Server:  &http.Server{},
		Logger:  log,
		Limiter: limiter,
	}, nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			))
	}

	return log
}
