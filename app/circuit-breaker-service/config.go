package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/server"
	yaml "gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	DefaultPageSize                int `yaml:"default_page_size"`
	DefaultErrorsThreshold         int `yaml:"default_errors_threshold"`
	DefaultResetTimeoutMs          int `yaml:"default_reset_timeout_ms"`
	DefaultErrorsCntResetTimeoutMs int `yaml:"default_errors_cnt_reset_timeout_ms"`
}

type Config struct {
	LogLevel string `yaml:"log_level"`

	API server.Config `yaml:"api"`
	// Database                   storage.Config                         `yaml:"db"`
	Service ServiceConfig `yaml:"service"`
}

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the YAML data into the Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (c *ServiceConfig) Validate() error {
	return nil
}

func (c *Config) Validate() error {
	if err := c.API.Validate(); err != nil {
		return fmt.Errorf("failed to validate HTTP Server Config, error: '%w'", err)
	}

	if err := c.Service.Validate(); err != nil {
		return fmt.Errorf("failed to validate service config, error: '%w'", err)
	}

	var logLevel = new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return fmt.Errorf("failed to validate log level, error: '%w'", err)
	}

	return nil
}

func configureLogging(cfg *Config) (*slog.Logger, error) {
	var logLevel = new(slog.LevelVar)
	// NOTE(maksym): the log_level param should already be validated
	_ = logLevel.UnmarshalText([]byte(cfg.LogLevel))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel, ReplaceAttr: hideAuthKey}))
	slog.SetDefault(logger)

	return logger, nil
}

func hideAuthKey(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "config" {
		config, ok := a.Value.Any().(Config)
		if ok {
			const asterisks = "***"
			config.API.AuthKey = asterisks
			a.Value = slog.AnyValue(config)
		}
	}

	return a
}
