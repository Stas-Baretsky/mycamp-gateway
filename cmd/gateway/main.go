package main

import (
	"context"
	"gateway-api/internal/bootstrap"
	"gateway-api/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	app, err := bootstrap.Build(ctx, *cfg)
	if err != nil {
		panic(err)
	}
	app.MustRun()

	// TODO : cache : redis
	// TODO : init router : chi
	// TODO : run server :
}
