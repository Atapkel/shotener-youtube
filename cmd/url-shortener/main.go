package main

import (
	config2 "github.com/atapkel/shortener/internal/config"
	"github.com/atapkel/shortener/internal/http-server/handlers/url/save"
	"github.com/atapkel/shortener/internal/http-server/middleware/logger"
	"github.com/atapkel/shortener/internal/lib/logger/sl"
	"github.com/atapkel/shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config2.MustLoad()
	log := setUpLogger(config.Env)
	strg, err := storage.New(config.StoragePath)
	if err != nil {
		log.Error("Failed to init sql", sl.Err(err))
		os.Exit(1)
	}
	strg.SaveURL("beibarys", "b")
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, strg))

	log.Info("Staring shortener server", slog.String("env", config.Env))
	log.Info("starting server", slog.String("address", config.HttpServer.Address))

	srv := &http.Server{
		Addr:         config.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  config.HttpServer.Timeout,
		WriteTimeout: config.HttpServer.Timeout,
		IdleTimeout:  config.HttpServer.IdleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("error on starting server")
	}
	log.Error("server stopped")
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
