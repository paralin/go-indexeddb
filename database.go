package indexeddb

import (
	"github.com/pkg/errors"
	"syscall/js"
)

// Database contains object stores, which contain data.
type Database struct {
	val js.Value
}

// NewDatabase constructs a database with a js object.
func NewDatabase(val js.Value) *Database {
	return &Database{val: val}
}

// GetName returns the database name.
func (d *Database) GetName() string {
	return d.val.Get("name").String()
}

// GetVersion returns the database version.
func (d *Database) GetVersion() int {
	return d.val.Get("version").Int()
}

// ContainsObjectStore checks if the db has a object store by id.
func (d *Database) ContainsObjectStore(id string) bool {
	return d.val.Get("objectStoreNames").Call("contains", id).Bool()
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
	switch mode {
	case READONLY:
	case READWRITE:
	default:
		mode = READONLY
	}

	scopeArg := make([]interface{}, len(scope))
	for i, x := range scope {
		scopeArg[i] = x
	}
	val := d.val.Call("transaction", scopeArg, string(mode))
	if !val.Truthy() {
		return nil, errors.Errorf("transaction(%v, %v): returned null", scope, mode)
	}

	return &Transaction{val: val}, nil
}

// GetJsValue returns the underlying js database handle.
func (d *Database) GetJsValue() js.Value {
	return d.val
}

// Close closes the database.
func (d *Database) Close() {
	d.val.Call("close")
}
