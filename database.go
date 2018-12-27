//+build js

package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// Database contains object stores, which contain data.
type Database struct {
	// Object is the database js object.
	*js.Object
}

// NewDatabase constructs a database with a js object.
func NewDatabase(obj *js.Object) *Database {
	return &Database{Object: obj}
}

// GetName returns the database name.
func (d *Database) GetName() string {
	return d.Get("name").String()
}

// GetVersion returns the database version.
func (d *Database) GetVersion() int {
	return d.Get("version").Int()
}

// ContainsObjectStore checks if the db has a object store by id.
func (d *Database) ContainsObjectStore(id string) bool {
	return d.Object.Get("objectStoreNames").Call("contains", id).Bool()
}

// TransactionMode is a transaction mode
type TransactionMode string

var (
	// READONLY is the read-only transaction mode
	READONLY TransactionMode = "readonly"
	// READWRITE is the read-write transaction mode
	READWRITE TransactionMode = "readwrite"
)

// Transaction gets a transaction with a object store name or names.
// Mode defaults to READONLY.
func (d *Database) Transaction(scope []string, mode TransactionMode) (t *Transaction, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()

	switch mode {
	case READONLY:
	case READWRITE:
	default:
		mode = READONLY
	}

	return &Transaction{Object: d.Object.Call("transaction", scope, mode)}, nil
}

// Close closes the database.
func (d *Database) Close() {
	d.Object.Call("close")
}
