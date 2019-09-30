package indexeddb

import (
	"errors"

	"syscall/js"
)

// DatabaseUpdate is a database during the updateneeded callback.
type DatabaseUpdate struct {
	*Database
}

// CreateObjectStoreOpts are the options for creating an object store.
type CreateObjectStoreOpts struct {
	val js.Value

	// keyPath is the key path.
	keyPath string
	// AutoIncrement if set
	autoIncrement bool
}

// NewCreateObjectStoreOpts constructs the options for CreateObjectStore.
func NewCreateObjectStoreOpts(keyPath string, autoIncrement bool) *CreateObjectStoreOpts {
	return &CreateObjectStoreOpts{
		keyPath:       keyPath,
		autoIncrement: autoIncrement,
	}
}

// ToJSValue converts the object to a js value.
func (o *CreateObjectStoreOpts) ToJSValue() js.Value {
	val := js.Global().Get("Object").New()
	val.Set("keyPath", o.keyPath)
	val.Set("autoIncrement", o.autoIncrement)
	return val
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
		args = append(args, opts.ToJSValue())
	}
	d.Database.val.Call("createObjectStore", args...)
	return nil
}
