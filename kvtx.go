// +build js

package indexeddb

import (
	"sync"
	"syscall/js"
)

// Kvtx implements a key-value transaction on top of DurableTransaction.
type Kvtx struct {
	txn         *DurableTransaction
	objStore    *DurableObjectStore
	discardOnce sync.Once
}

// NewKvtxTx constructs a new tranasction, opening the object store.
func NewKvtxTx(txn *DurableTransaction, objStoreID string) (*Kvtx, error) {
	objStore, err := txn.GetObjectStore(objStoreID)
	if err != nil {
		return nil, err
	}

	return &Kvtx{
		txn:      txn,
		objStore: objStore,
	}, nil
}

// Size returns the number of keys in the store.
func (t *Kvtx) Size() (uint64, error) {
	c, err := t.objStore.Count(nil)
	return uint64(c), err
}

// Get returns values for a key.
func (t *Kvtx) Get(key []byte) (data []byte, found bool, err error) {
	if len(key) == 0 {
		return nil, false, ErrEmptyKey
	}
	jsObj, err := t.objStore.Get(key)
	if err != nil {
		return nil, false, err
	}
	if !jsObj.Truthy() {
		return nil, false, nil
	}
	dlen := jsObj.Length()
	data = make([]byte, dlen)
	js.CopyBytesToGo(data, jsObj)
	return data, true, nil
}

// Set sets the value of a key.
// This will not be committed until Commit is called.
func (t *Kvtx) Set(key, value []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	return t.objStore.Put(value, key)
}

// Delete deletes a key.
// This will not be committed until Commit is called.
// Not found should not return an error.
func (t *Kvtx) Delete(key []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	return t.objStore.Delete(key)
}

// scanPrefix iterates over items with a prefix.
func (t *Kvtx) scanPrefix(prefix []byte, cb func(v *CursorValue) error) error {
	krv := js.Undefined()
	if len(prefix) != 0 {
		prefixGreater := make([]byte, len(prefix)+1)
		copy(prefixGreater, prefix)
		prefixGreater[len(prefixGreater)-1] = ^byte(0)
		krv = Bound(prefix, prefixGreater, false, false)
	}
	cursor, err := t.objStore.OpenCursor(krv)
	if err != nil {
		return err
	}
	for {
		val := cursor.WaitValue()
		if val == nil {
			return nil
		}

		if err := cb(val); err != nil {
			return err
		}

		cursor.ContinueCursor()
	}
}

// ScanPrefixKeys iterates over keys with a prefix.
func (t *Kvtx) ScanPrefixKeys(prefix []byte, cb func(key []byte) error) error {
	return t.scanPrefix(prefix, func(val *CursorValue) error {
		return cb(
			CopyByteSliceFromJs(val.Key),
		)
	})
}

// ScanPrefix iterates over keys with a prefix.
func (t *Kvtx) ScanPrefix(prefix []byte, cb func(key, val []byte) error) error {
	return t.scanPrefix(prefix, func(val *CursorValue) error {
		return cb(
			CopyByteSliceFromJs(val.Key),
			CopyByteSliceFromJs(val.Value),
		)
	})
}

// Exists checks if a key exists.
func (t *Kvtx) Exists(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, ErrEmptyKey
	}
	i, err := t.objStore.Count(key)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

// Commit commits the transaction to storage.
// Can return an error to indicate tx failure.
func (t *Kvtx) Commit() error {
	return t.txn.Commit()
}

// Discard cancels the transaction.
// If called after Commit, does nothing.
// Cannot return an error.
// Can be called unlimited times.
func (t *Kvtx) Discard() {
	t.txn.Abort()
}
