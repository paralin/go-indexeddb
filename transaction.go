package indexeddb

import (
	"github.com/pkg/errors"
	"syscall/js"
)

// Transaction is a database transaction.
type Transaction struct {
	val js.Value
}

// GetMode returns the transaction mode.
func (t *Transaction) GetMode() TransactionMode {
	return TransactionMode(t.val.Get("mode").String())
}

// GetObjectStore returns a object store.
func (t *Transaction) GetObjectStore(id string) (o *ObjectStore, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()

	s := t.val.Call("objectStore", id)
	if !s.Truthy() {
		return nil, errors.Errorf("GetObjectStore(%s) returned nil", id)
	}
	return &ObjectStore{val: s}, nil
}

// Abort aborts a transaction.
func (t *Transaction) Abort() {
	t.val.Call("abort")
}
