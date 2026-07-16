package main

import (
	"fmt"
	"gateway-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO : init logger : slog
	// TODO : init storage : postgresql
	// TODO : init router : chi
	// TODO : run server :
}
