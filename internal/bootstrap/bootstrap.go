package bootstrap

import (
	"gateway-api/internal/config"
	"gateway-api/internal/app"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Build(cfg config.Config) (app.Application, error) {

	log := setupLogger(cfg.Env)
	log.Info("Starting service...", slog.String("env", cfg.Env))

	return app.Application{
		Server: ,
		Logger: log,
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("Env must be prod, dev or local")
	}
	return log
}
