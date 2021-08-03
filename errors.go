package indexeddb

import (
	"errors"
	"strings"
)

var (
	// ErrEmptyKey is returned if the key was empty.
	ErrEmptyKey = errors.New("key cannot be empty")
)

// errIsInactiveTransaction checks if an error is the "inactive transaction" error
func errIsInactiveTransaction(err error) bool {
	return strings.Contains(err.Error(), "transaction is not active")
}
