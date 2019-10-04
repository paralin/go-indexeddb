// +build js,wasm

package indexeddb

import (
	"context"
	"errors"

	"syscall/js"
)

// IndexedDB is the global indexeddb object.
type IndexedDB struct {
	val js.Value
}

// GlobalIndexedDB returns the global IndexedDB object.
func GlobalIndexedDB() *IndexedDB {
	val := js.Global().Get("indexedDB")
	if !val.Truthy() {
		return nil
	}
	return &IndexedDB{val: val}
}

// Open opens an indexeddb database with a version and upgrader.
func (i *IndexedDB) Open(
	ctx context.Context,
	name string,
	version int,
	upgrader func(d *DatabaseUpdate, oldVersion, newVersion int) error,
) (*Database, error) {
	var db *Database
	errCh := make(chan error, 1)
	putErr := func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}
	odbReq := i.val.Call("open", name, version)
	odbReq.Set("onupgradeneeded", js.FuncOf(
		func(th js.Value, dats []js.Value) interface{} {
			event := dats[0]
			// event is an IDBVersionChangeEvent
			oldVersion := event.Get("oldVersion").Int()
			newVersion := event.Get("newVersion").Int()
			db = &Database{val: event.Get("target").Get("result")}
			if err := upgrader(&DatabaseUpdate{Database: db}, oldVersion, newVersion); err != nil {
				putErr(err)
			}
			return nil
		},
	))
	odbReq.Set("onerror", js.FuncOf(
		func(th js.Value, dats []js.Value) interface{} {
			o := dats[0]
			go putErr(errors.New(o.
				Get("target").
				Get("error").
				Get("message").
				String(),
			))
			return nil
		},
	))
	odbReq.Set("onsuccess", js.FuncOf(
		func(th js.Value, dats []js.Value) interface{} {
			o := dats[0]
			if db == nil {
				db = NewDatabase(o.Get("target").Get("result"))
			}
			go putErr(nil)
			return nil
		},
	))
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// GetJsValue returns the underlying js database handle.
func (i *IndexedDB) GetJsValue() js.Value {
	return i.val
}
