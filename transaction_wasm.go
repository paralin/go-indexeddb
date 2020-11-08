// +build js,wasm

package indexeddb

import (
	"sync"
	"syscall/js"

	"github.com/pkg/errors"
)

// Transaction is a database transaction.
type Transaction struct {
	val       js.Value
	abortOnce sync.Once
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

// GetJsValue returns the underlying js database handle.
func (t *Transaction) GetJsValue() js.Value {
	return t.val
}

// Abort aborts a transaction.
func (t *Transaction) Abort() {
	defer func() {
		if err := recover(); err != nil {
			// ignore error here
			// Failed to execute 'abort' on 'IDBTransaction': The transaction has finished.
			_ = err
		}
	}()
	t.abortOnce.Do(func() {
		t.val.Call("abort")
	})
}
