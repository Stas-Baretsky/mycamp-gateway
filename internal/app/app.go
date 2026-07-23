package app

import (
	"context"
	"gateway-api/internal/bootstrap"
	"gateway-api/internal/config"
	"gateway-api/internal/middleware/ratelimiter"
	"log/slog"
	"net/http"
)

type Application struct {
	Server  *http.Server
	Logger  *slog.Logger
	Limiter *ratelimiter.RateLimiter
}

func (a *Application) MustRun(ctx context.Context, cfg config.Config) {
	app, err := bootstrap.Build(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return &app
}

func (a *Application) Stop() {
	///server stop
	///logger stop
	///limiter stop
}
