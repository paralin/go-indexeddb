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

// Open opens an indexeddb database.
func (i *IndexedDB) Open(ctx context.Context, name string) (*Database, error) {
	odbReq := i.Call("open", name)
	errCh := make(chan error, 1)
	putErr := func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}
	odbReq.Set("onerror", func(o *js.Object) {
		go putErr(errors.New("cannot open indexeddb"))
	})
	odbReq.Set("onsuccess", func(o *js.Object) {
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

	return &Database{Object: odbReq.Get("result")}, nil
}

// OpenWithUpgrader opens an indexeddb database with a version and upgrader.
/*
func (i *IndexedDB) OpenWithUpgrader(
	ctx context.Context,
	name string,
	ver int,
	upgrader func(TODO),
) (*Database, error) {
	odbReq := i.Call("open", name)
}
*/
