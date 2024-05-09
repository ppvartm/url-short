package main

import (
	"fmt"
	"log/slog"
	"urlshort/internal/config"
	"urlshort/pkg/logger"
)

func main() {
	config := config.MustLoad()

	fmt.Println(*config)

	log := logger.SetupLogger(config.Env)

	log.Debug("start", slog.String("env", config.Env))

	// TODO: init storage

	// TODO: init router

	// TODO: run server

}
