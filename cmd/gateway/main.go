package main

import (
	"gateway-api/internal/bootstrap"
	"gateway-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app, err := bootstrap.Build(*cfg)

	// TODO : cache : reddis
	// TODO : init router : chi
	// TODO : run server :
}
