package map_test_storage

import (
	"context"
	"errors"
	"sync"

	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/generic_storage"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/model"
)

// Shutdown gracefully shuts down the storage client.
func (c *Client) Shutdown(ctx context.Context) error {
	if !c.initialized {
		return generic_storage.ErrNotInitialized
	}
	c.logger.Debug("Shutdown called")
	c.registry = sync.Map{} // Clear all entries
	return nil
}

// IsAlive checks if the storage client is alive.
func (c *Client) IsAlive(ctx context.Context) error {
	if !c.initialized {
		return generic_storage.ErrNotInitialized
	}
	c.logger.Debug("IsAlive called")
	return nil
}

// UpsertEntry inserts or updates an entry in the storage.
func (c *Client) UpsertEntry(ctx context.Context, primaryKey model.Key, entry model.CircuitBreakerEntry) error {
	if !c.initialized {
		return generic_storage.ErrNotInitialized
	}
	c.logger.Debug("UpsertEntry called", "primaryKey", primaryKey, "entry", entry)
	c.registry.Store(primaryKey, entry)
	return nil
}

// AddNewEntry adds a new entry to the storage. Fails if the key already exists.
func (c *Client) AddNewEntry(ctx context.Context, primaryKey model.Key, entry model.CircuitBreakerEntry) error {
	if !c.initialized {
		return generic_storage.ErrNotInitialized
	}
	c.logger.Debug("AddNewEntry called", "primaryKey", primaryKey, "entry", entry)
	if _, exists := c.registry.Load(primaryKey); exists {
		err := errors.New("entry already exists")
		c.logger.Debug("AddNewEntry failed", "primaryKey", primaryKey, "error", err)
		return err
	}
	c.registry.Store(primaryKey, entry)
	return nil
}

func (c *Client) RemoveEntry(ctx context.Context, primaryKey model.Key) error {
	if !c.initialized {
		return generic_storage.ErrNotInitialized
	}
	c.logger.Debug("RemoveEntry called", "primaryKey", primaryKey)
	if _, exists := c.registry.Load(primaryKey); !exists {
		err := errors.New("entry does not exist")
		c.logger.Debug("RemoveEntry failed", "primaryKey", primaryKey, "error", err)
		return err
	}
	c.registry.Delete(primaryKey)
	return nil
}

func (c *Client) GetEntry(ctx context.Context, primaryKey model.Key) (model.CircuitBreakerEntry, error) {
	if !c.initialized {
		return model.CircuitBreakerEntry{}, generic_storage.ErrNotInitialized
	}
	c.logger.Debug("GetEntry called", "primaryKey", primaryKey)
	entry, exists := c.registry.Load(primaryKey)
	if !exists {
		err := errors.New("entry does not exist")
		c.logger.Debug("GetEntry failed", "primaryKey", primaryKey, "error", err)
		return model.CircuitBreakerEntry{}, err
	}
	return entry.(model.CircuitBreakerEntry), nil
}

// GetAllEntries retrieves all entries in the storage.
func (c *Client) GetAllEntries(ctx context.Context) ([]model.CircuitBreakerEntry, error) {
	if !c.initialized {
		return nil, generic_storage.ErrNotInitialized
	}
	c.logger.Debug("GetAllEntries called")
	var entries []model.CircuitBreakerEntry
	c.registry.Range(func(_, value interface{}) bool {
		entries = append(entries, value.(model.CircuitBreakerEntry))
		return true
	})
	return entries, nil
}

// GetAllEntriesPaginated retrieves a paginated list of entries starting from the lastPrimaryKey.
func (c *Client) GetAllEntriesPaginated(ctx context.Context, lastPrimaryKey model.Key, pageSize int) ([]model.CircuitBreakerEntry, error) {
	if !c.initialized {
		return nil, generic_storage.ErrNotInitialized
	}
	c.logger.Debug("GetAllEntriesPaginated called", "lastPrimaryKey", lastPrimaryKey, "pageSize", pageSize)
	var entries []model.CircuitBreakerEntry
	var foundLastKey bool

	c.registry.Range(func(key, value interface{}) bool {
		if foundLastKey || key.(model.Key) == lastPrimaryKey {
			foundLastKey = true
			entries = append(entries, value.(model.CircuitBreakerEntry))
			if len(entries) >= pageSize {
				return false
			}
		}
		return true
	})

	return entries, nil
}

// GetAllPrimaryKeys retrieves all primary keys in the storage.
func (c *Client) GetAllPrimaryKeys(ctx context.Context) ([]model.Key, error) {
	if !c.initialized {
		return nil, generic_storage.ErrNotInitialized
	}
	c.logger.Debug("GetAllPrimaryKeys called")
	var keys []model.Key
	c.registry.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(model.Key))
		return true
	})
	return keys, nil
}
