package app

import (
	"gateway-api/internal/middleware/ratelimiter"
	"log/slog"
	"net/http"
)

type Application struct {
	Server  *http.Server
	Logger  *slog.Logger
	Limiter *ratelimiter.RateLimiter
}

func (a *Application) MustRun() {
	///
}

func (a *Application) Stop() {
	///
}
