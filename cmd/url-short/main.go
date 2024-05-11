package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"urlshort/internal/config"
	"urlshort/internal/http-server/handlers/redirect"
	"urlshort/internal/http-server/handlers/url/save"
	mwLogger "urlshort/internal/http-server/middleware/logger"
	"urlshort/internal/lib/logger"
	"urlshort/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("server starting", slog.String("addres", config.Address))

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}
