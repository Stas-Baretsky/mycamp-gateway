package bootstrap

import (
	"context"
	"fmt"

	"gateway-api/internal/app"
	limiterCache "gateway-api/internal/cache/limiterCache"
	"gateway-api/internal/client/auth"
	"gateway-api/internal/config"
	"gateway-api/internal/middleware/ratelimiter"
	"log/slog"
	"net/http"
)

func Build(ctx context.Context, cfg config.Config, log slog.Logger) (*app.Application, error) {
	const op = "bootstrap.bootstrap.Build"

	log.Info("Starting service...", slog.String("env", cfg.Env))

	redisStorage, err := limiterCache.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	limiter, err := ratelimiter.NewRateLimiter(ctx, &log, redisStorage, &cfg.RateLimiter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	authClient, err := auth.New(&log, cfg)

	return &app.Application{
		Server: &http.Server{
			Addr:        cfg.HTTPServer.Address,
			IdleTimeout: cfg.HTTPServer.IddleTimeout,
			ReadTimeout: cfg.HTTPServer.Timeout,
		},
		AuthClient: authClient,
		Log:        &log,
		Limiter:    limiter,
	}, nil
}
