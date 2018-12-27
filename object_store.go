//+build js

package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// ObjectStore is a object store attached to a transaction.
type ObjectStore struct {
	// Object is the database js object.
	*js.Object
}

// GetName returns the object store name.
func (s *ObjectStore) GetName() string {
	return s.Object.Get("name").String()
}

// Put puts data into the store.
func (s *ObjectStore) Put(value interface{}, optionalKey interface{}) error {
	_, err := WaitRequest(s.Object.Call("put", value, optionalKey))
	return err
}

// Add adds data to the store.
func (s *ObjectStore) Add(value interface{}, optionalKey interface{}) error {
	_, err := WaitRequest(s.Object.Call("add", value, optionalKey))
	return err
}

// Delete deletes data from the store.
func (s *ObjectStore) Delete(query interface{}) error {
	_, err := WaitRequest(s.Object.Call("delete", query))
	return err
}

// Clear clears all data from the store.
func (s *ObjectStore) Clear() error {
	_, err := WaitRequest(s.Object.Call("clear"))
	return err
}

// Get gets data from the store
func (s *ObjectStore) Get(query interface{}) (*js.Object, error) {
	return WaitRequest(s.Object.Call("get", query))
}

// GetKey gets data from the store by key.
func (s *ObjectStore) GetKey(query interface{}) (*js.Object, error) {
	return WaitRequest(s.Object.Call("getKey", query))
}

// GetAll gets all values matching an optional query with an optional count.
func (s *ObjectStore) GetAll(query interface{}) (*js.Object, error) {
	return WaitRequest(s.Object.Call("getAll", query))
}

// GetAllKeys gets all keys matching an optional query with an optional count.
func (s *ObjectStore) GetAllKeys(query interface{}) (*js.Object, error) {
	return WaitRequest(s.Object.Call("getAllKeys", query))
}

// Count counts keys matching the optional query.
func (s *ObjectStore) Count(query interface{}) (*js.Object, error) {
	return WaitRequest(s.Object.Call("count", query))
}
