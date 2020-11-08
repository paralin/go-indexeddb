// +build js,!wasm

package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// Transaction is a database transaction.
type Transaction struct {
	// Object is the database js object.
	*js.Object
}

// WaitTransactionComplete waits for oncomplete on a transaction.
// Registers onsuccess and onerror.
// Returns transaction.error if set, or nil.
// Call commit before calling this.
func WaitTransactionComplete(obj *js.Object) error {
	ret := func() error {
		var err error
		if o := obj.Get("error"); o != nil && o != js.Undefined {
			err = errors.New(o.Get("message").String())
		}
		return err
	}
	errCh := make(chan struct{}, 1)
	rerr := func() {
		select {
		case errCh <- struct{}{}:
		default:
		}
	}
	obj.Set("onerror", func(e *js.Object) {
		rerr()
	})
	obj.Set("oncomplete", func(e *js.Object) {
		rerr()
	})
	<-errCh
	return ret()
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

// Abort aborts a transaction, rolling back changes.
//
// Note: transactions auto-commit.
func (t *Transaction) Abort() {
	t.Object.Call("abort")
}

// Commit forces committing a transaction.
//
// Note: transactions auto-commit, abort will roll-back the changes.
func (t *Transaction) Commit() {
	t.Object.Call("commit")
}

// WaitComplete waits for the transaction to complete.
// Call commit() first.
// Returns any error if set on the transaction.
func (t *Transaction) WaitComplete() error {
	return WaitTransactionComplete(t.Object)
}
