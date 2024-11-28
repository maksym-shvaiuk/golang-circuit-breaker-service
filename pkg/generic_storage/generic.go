package generic_storage

import (
	"context"
)

// NOTE: K for primary key, T for full entry.
type StorageClient[K any, T any] interface {
	Shutdown(ctx context.Context) error
	IsAlive(ctx context.Context) error

	// NOTE (maksym): updates attribute, creates attribute if it not exists
	UpsertEntry(ctx context.Context, primaryKey K, entry T) error
	AddNewEntry(ctx context.Context, primaryKey K, entry T) error
	RemoveEntry(ctx context.Context, primaryKey K) error
	GetEntry(ctx context.Context, primaryKey K) (T, error)
	GetAllEntries(ctx context.Context) ([]T, error)
	GetAllEntriesPaginated(ctx context.Context, lastPrimaryKey K, pageSize int) ([]T, error)
	GetAllPrimaryKeys(ctx context.Context) ([]K, error)

	// TODO: implement RemoveEntriesbatch(), if necessary to not establish DB connection on each entry
	// TODO: implement AddNewEntriesbatch(), if necessary to not establish DB connection on each entry
	// TODO: implement UpsertEntriesbatch(), if necessary to not establish DB connection on each entry
}
