//+build js

package indexeddb

import (
	"context"
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// IndexedDB is the global indexeddb object.
type IndexedDB struct {
	*js.Object
}

// GlobalIndexedDB returns the global IndexedDB object.
func GlobalIndexedDB() *IndexedDB {
	return &IndexedDB{Object: js.Global.Get("indexedDB")}
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
	odbReq := i.Call("open", name, version)
	odbReq.Set("onupgradeneeded", func(event *js.Object) {
		// event is an IDBVersionChangeEvent
		oldVersion := event.Get("oldVersion").Int()
		newVersion := event.Get("newVersion").Int()
		db = &Database{Object: event.Get("target").Get("result")}
		if err := upgrader(&DatabaseUpdate{Database: db}, oldVersion, newVersion); err != nil {
			putErr(err)
		}
	})
	odbReq.Set("onerror", func(o *js.Object) {
		js.Global.Set("idbError", o)
		go putErr(errors.New(o.
			Get("target").
			Get("error").
			Get("message").
			String(),
		))
	})
	odbReq.Set("onsuccess", func(o *js.Object) {
		if db == nil {
			db = &Database{Object: o.Get("target").Get("result")}
		}
		go putErr(nil)
	})
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
