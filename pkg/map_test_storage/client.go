package map_test_storage

import (
	"log/slog"
	"sync"
)

// NOTE (maksym): this dummy storage should be used in unit tests only

var client Client

type Client struct {
	logger      *slog.Logger
	registry    sync.Map
	initialized bool
}

func New(logger *slog.Logger) (*Client, error) {
	if !client.initialized {
		client = Client{
			logger:      logger.With("component", "test-storage"),
			registry:    sync.Map{},
			initialized: true,
		}
	}

	return &client, nil
}
