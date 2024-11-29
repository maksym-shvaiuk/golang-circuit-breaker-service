package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/generic_storage"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/model"
)

type Service struct {
	storage generic_storage.StorageClient[model.Key, model.CircuitBreakerEntry]
	logger  *slog.Logger
	engine  *gin.Engine
}

type Config struct {
	ServerHost       string        `yaml:"server_host"`
	ServerPort       int           `yaml:"server_port"`
	RequestRWTimeout time.Duration `yaml:"request_rw_timeout"`
	IdleTimeout      time.Duration `yaml:"idle_timeout"`
	GracefulTimeout  time.Duration `yaml:"graceful_timeout"`
	AuthKey          string        `yaml:"auth_key"`
}

func (c *Config) Validate() error {
	if c.ServerHost == "" {
		return fmt.Errorf("missed ServerHost config param")
	}

	if c.ServerPort == 0 {
		return fmt.Errorf("missed ServerPort config param")
	}

	if c.AuthKey == "" {
		return fmt.Errorf("missed AuthKey config param")
	}

	return nil
}
