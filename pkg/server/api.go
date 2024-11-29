package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/generic_storage"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/model"
)

func New(cfg *Config, storage generic_storage.StorageClient[model.Key, model.CircuitBreakerEntry], logger *slog.Logger) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	if storage == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Initialize Gin engine
	engine := gin.Default()

	// Add middlewares
	engine.Use(loggingMiddleware(logger))
	engine.Use(authMiddlewareWithToken(cfg.AuthKey))

	// Register routes
	registerRoutes(engine)

	return &Service{
		storage: storage,
		logger:  logger,
		engine:  engine,
	}, nil
}

func (s *Service) Run(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if s.engine == nil {
		return fmt.Errorf("gin engine is not initialized")
	}

	address := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	s.logger.Info("Starting server", "address", address)

	server := &http.Server{
		Addr:         address,
		Handler:      s.engine,
		ReadTimeout:  cfg.RequestRWTimeout,
		WriteTimeout: cfg.RequestRWTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		<-context.Background().Done()
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error("Server shutdown failed", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
