package bootstrap

import (
	"context"
	"gateway-api/internal/app"
	goredis "gateway-api/internal/cache/go-redis"
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
	log := setupLogger(cfg.Env)
	log.Info("Starting service...", slog.String("env", cfg.Env))

	redisStorage, err := goredis.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	limiter, err := ratelimiter.NewRateLimiter(ctx, log, redisStorage, &cfg.RateLimiter)

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
