package main_test

import (
	"testing"

	main "github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/app/circuit-breaker-service"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/server"
)

func TestValidateConfigValid(t *testing.T) {
	// Arrange
	validConfig := main.Config{
		LogLevel: "info",
		API: server.Config{
			ServerHost: "localhost",
			ServerPort: 8080,
			AuthKey:    "valid-auth-key",
		},
		Service: main.ServiceConfig{
			DefaultPageSize:                5,
			DefaultErrorsThreshold:         10,
			DefaultResetTimeoutMs:          60000,
			DefaultErrorsCntResetTimeoutMs: 10000,
		},
	}

	// Act
	err := validConfig.Validate()

	// Assert
	if err != nil {
		t.Fatalf("Expected no validation errors, but got: %v", err)
	}
}

func TestValidateConfigInvalidHttpServer(t *testing.T) {
	// Arrange
	invalidConfig := main.Config{
		LogLevel: "info",
		API:      server.Config{},
		Service: main.ServiceConfig{
			DefaultPageSize:                5,
			DefaultErrorsThreshold:         10,
			DefaultResetTimeoutMs:          60000,
			DefaultErrorsCntResetTimeoutMs: 10000,
		},
	}

	// Act
	err := invalidConfig.Validate()

	// Assert
	if err == nil {
		t.Fatalf("Expected validation error for invalid HTTP Server config, but got none")
	}
}

func TestValidateConfigInvalidLogLevel(t *testing.T) {
	// Arrange
	invalidConfig := main.Config{
		LogLevel: "invalid",
		API: server.Config{
			ServerHost: "localhost",
			ServerPort: 8080,
			AuthKey:    "valid-auth-key",
		},
		Service: main.ServiceConfig{
			DefaultPageSize:                5,
			DefaultErrorsThreshold:         10,
			DefaultResetTimeoutMs:          60000,
			DefaultErrorsCntResetTimeoutMs: 10000,
		},
	}

	// Act
	err := invalidConfig.Validate()

	// Assert
	if err == nil {
		t.Fatalf("Expected validation error for invalid log level, but got none")
	}
}
