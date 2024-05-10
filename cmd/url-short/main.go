package main

import (
	"fmt"
	"log/slog"
	"os"
	"urlshort/internal/config"
	"urlshort/internal/lib/logger"
	"urlshort/internal/storage/sqlite"
)

func main() {
	config := config.MustLoad()

	fmt.Println(*config)

	log := logger.SetupLogger(config.Env)

	log.Debug("start", slog.String("env", config.Env))

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router

	// TODO: run server

}
