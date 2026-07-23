package app

import (
	"context"
	"gateway-api/internal/config"
	"gateway-api/internal/middleware/ratelimiter"
	"log/slog"
	"net/http"
)

type Application struct {
	Server  *http.Server
	Log     *slog.Logger
	Limiter *ratelimiter.RateLimiter
}

func (a *Application) Run(ctx context.Context, cfg config.Config) error {
	a.Log.Info("Starting http server...")
	return a.Server.ListenAndServe()
	//
}

func (a *Application) Stop(ctx context.Context) error {
	a.Log.Info("Stopping application...")

	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	if a.Limiter != nil {
		a.Limiter.Stop(ctx)
	}

	return nil
}
