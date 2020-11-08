package indexeddb

import "strings"

// errIsInactiveTransaction checks if an error is the "inactive transaction" error
func errIsInactiveTransaction(err error) bool {
	return strings.Contains(err.Error(), "transaction is not active")
}
