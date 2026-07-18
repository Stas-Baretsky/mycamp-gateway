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

	// TODO : cache : reddis
	// TODO : init router : chi
	// TODO : run server :
}
