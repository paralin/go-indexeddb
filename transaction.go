package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// Transaction is a database transaction.
type Transaction struct {
	// Object is the database js object.
	*js.Object
}

// GetMode returns the transaction mode.
func (t *Transaction) GetMode() TransactionMode {
	return TransactionMode(t.Object.Get("mode").String())
}

// GetObjectStore returns a object store.
func (t *Transaction) GetObjectStore(id string) (o *ObjectStore, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()

	s := t.Object.Call("objectStore", id)
	return &ObjectStore{Object: s}, nil
}

// Abort aborts a transaction.
func (t *Transaction) Abort() {
	t.Object.Call("abort")
}
