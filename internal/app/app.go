package app

import (
	"log/slog"
	"net/http"
)

type Application struct {
	Server *http.Server
	Logger *slog.Logger
}

func (a *Application) MustRun() {
	///
}

func (a *Application) Stop() {
	///
}
