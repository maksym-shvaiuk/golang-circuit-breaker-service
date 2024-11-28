package generic_storage

import (
	"errors"
)

var ErrNotInitialized = errors.New("storage: not initialized")
var ErrEntryNotFound = errors.New("storage: entry not found")
var ErrEntryAlreadyExists = errors.New("storage: entry already exists")
