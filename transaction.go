// +build js

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

// WaitTransactionComplete waits for oncomplete on a transaction.
// Registers onsuccess and onerror.
// Returns transaction.error if set, or nil.
// Call commit before calling this.
func WaitTransactionComplete(obj js.Value) error {
	ret := func() error {
		var err error
		if o := obj.Get("error"); o.Truthy() {
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
	cb := js.FuncOf(func(th js.Value, dats []js.Value) interface{} {
		go rerr()
		return nil
	})
	obj.Set("onerror", cb)
	obj.Set("oncomplete", cb)
	<-errCh
	return ret()
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

// Abort aborts a transaction, rolling back changes.
//
// Note: transactions auto-commit.
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

// Commit forces committing a transaction.
//
// Note: transactions auto-commit, abort will roll-back the changes.
func (t *Transaction) Commit() {
	defer func() {
		if err := recover(); err != nil {
			// ignore error here
			_ = err
		}
	}()
	t.abortOnce.Do(func() {
		t.val.Call("commit")
	})
}

// WaitComplete waits for the transaction to complete.
// Call commit() first.
// Returns any error if set on the transaction.
func (t *Transaction) WaitComplete() error {
	return WaitTransactionComplete(t.val)
}
