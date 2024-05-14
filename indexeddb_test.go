//go:build js
// +build js

package indexeddb

import (
	"context"
	"errors"
	"testing"
)

func TestIndexedDB(t *testing.T) {
	ctx := context.Background()
	ver := 1
	id := "testObjectStore"
	
	// Open database
	db, err := GlobalIndexedDB().Open(
		ctx,
		"test-db",
		ver,
		func(d *DatabaseUpdate, oldVersion, newVersion int) error {
			if !d.ContainsObjectStore(id) {
				if err := d.CreateObjectStore(id, nil); err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	
	// Get durable transaction
	durTx, err := NewDurableTransaction(db, []string{id}, READWRITE)
	if err != nil {
		t.Fatalf("Error getting durable transaction: %v", err)
	}

	// Get object store
	objStore, err := NewKvtxTx(durTx, id)
	if err != nil {
		t.Fatalf("Error getting object store: %v", err)
	}

	// Set key/value
	key := []byte("key")
	val := []byte("test")
	if err := objStore.Set(key, val); err != nil {
		t.Fatalf("Error setting key/value: %v", err)
	}

	if err := objStore.Commit(); err != nil {
		t.Fatalf("Error committing transaction: %v", err)
	}

	// Get value
	dat, found, err := objStore.Get(key)
	if err == nil && !found {
		err = errors.New("key not found after setting it")
	}
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}
	if string(dat) != string(val) {
		t.Fatalf("Got wrong value. Expected %s, got %s", val, dat)
	}

	// Scan prefix
	prefix := []byte("ke")
	err = objStore.ScanPrefix(prefix, func(key, val []byte) error {
		if string(key) != string(key) {
			t.Errorf("Wrong key returned. Expected %s, got %s", key, key)
		}
		if string(val) != string(val) {
			t.Errorf("Wrong value returned. Expected %s, got %s", val, val)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error scanning prefix: %v", err)
	}
}
