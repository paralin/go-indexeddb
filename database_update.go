//+build js

package indexeddb

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// DatabaseUpdate is a database during the updateneeded callback.
type DatabaseUpdate struct {
	*Database
}

// CreateObjectStoreOpts are the options for creating an object store.
type CreateObjectStoreOpts struct {
	*js.Object

	// KeyPath is the key path.
	KeyPath string `js:"keyPath"`
	// AutoIncrement if set
	AutoIncrement bool `js:"autoIncrement"`
}

// NewCreateObjectStoreOpts constructs the options for CreateObjectStore.
func NewCreateObjectStoreOpts(keyPath string, autoIncrement bool) *CreateObjectStoreOpts {
	o := &CreateObjectStoreOpts{
		Object: js.Global.Get("Object").New(),
	}
	o.KeyPath = keyPath
	o.AutoIncrement = autoIncrement
	return o
}

// CreateObjectStore creates an object store.
// keyPath is optional
func (d *DatabaseUpdate) CreateObjectStore(
	id string,
	opts *CreateObjectStoreOpts,
) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			if err == nil {
				var ok bool
				err, ok = rerr.(error)
				if !ok {
					err = errors.New("create object store paniced")
				}
			}
		}
	}()

	args := []interface{}{id}
	if opts != nil {
		args = append(args, opts)
	}
	d.Object.Call("createObjectStore", args...)
	return nil
}
