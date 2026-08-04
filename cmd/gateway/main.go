package main

import (
	"context"
	"gateway-api/internal/bootstrap"
	"gateway-api/internal/config"
	"gateway-api/internal/handlers/url/register"
	"gateway-api/internal/lib/logger/sl"
	"gateway-api/internal/middleware/mwlogger"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	app, err := bootstrap.Build(ctx, *cfg, *log)
	if err != nil {
		log.Error("Failed to building application", sl.Err(err))
		panic(err)
	}
	app.Run(ctx, *cfg)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post(
		"/auth/register",
		register.New(log, app.AuthClient))

	// TODO : cache : redis
	// TODO : init router : chi
	// TODO : run server :
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
