//+build js

package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// Database contains object stores, which contain data.
type Database struct {
	// Object is the database js object.
	*js.Object

	// name is the database name
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

// Close closes the database.
func (d *Database) Close() {
	d.Object.Call("close")
}
