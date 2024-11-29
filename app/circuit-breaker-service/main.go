package main

import (
	"os"

	"log/slog"

	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/map_test_storage"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/server"
)

func main() {
	initFlags()

	cfg, err := loadConfig(*FlagConfigFilePath)
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger, err := configureLogging(cfg)
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		logger.Error("failed to configure logging", "error", err)
		os.Exit(2)
	}

	logger.Info("starting service")
	logger.Info("config loaded", "config", *cfg)

	storage, _ := map_test_storage.New(logger)

	service, err := server.New(&cfg.API, storage, logger)
	if err != nil {
		logger.Error("Failed to initialize service", "error", err)
		os.Exit(3)
	}

	if err := service.Run(&cfg.API); err != nil {
		logger.Error("Server encountered an error", "error", err)
		os.Exit(1)
	}
}
