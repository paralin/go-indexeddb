// +build js

package indexeddb

import (
	"syscall/js"
)

// ObjectStore is a object store attached to a transaction.
type ObjectStore struct {
	val js.Value
}

// GetName returns the object store name.
func (s *ObjectStore) GetName() string {
	return s.val.Get("name").String()
}

// Put puts data into the store.
func (s *ObjectStore) Put(value interface{}, key interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	value = MaybeConvertValueToJs(value)
	key = MaybeConvertValueToJs(key)
	_, err := WaitRequest(s.val.Call("put", value, key))
	return err
}

// Add adds data to the store.
func (s *ObjectStore) Add(value interface{}, key interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	value = MaybeConvertValueToJs(value)
	key = MaybeConvertValueToJs(key)
	_, err := WaitRequest(s.val.Call("add", value, key))
	return err
}

// Delete deletes data from the store.
func (s *ObjectStore) Delete(query interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	query = MaybeConvertValueToJs(query)
	_, err := WaitRequest(s.val.Call("delete", query))
	return err
}

// Clear clears all data from the store.
func (s *ObjectStore) Clear() (e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	_, err := WaitRequest(s.val.Call("clear"))
	return err
}

// Get gets data from the store
func (s *ObjectStore) Get(query interface{}) (rv js.Value, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	query = MaybeConvertValueToJs(query)
	return WaitRequest(s.val.Call("get", query))
}

// GetKey gets data from the store by key.
func (s *ObjectStore) GetKey(query interface{}) (_ js.Value, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	query = MaybeConvertValueToJs(query)
	return WaitRequest(s.val.Call("getKey", query))
}

// GetAll gets all values matching an optional query with an optional count.
func (s *ObjectStore) GetAll(query interface{}) (_ js.Value, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	query = MaybeConvertValueToJs(query)
	return WaitRequest(s.val.Call("getAll", query))
}

// GetAllKeys gets all keys matching an optional query with an optional count.
func (s *ObjectStore) GetAllKeys(query interface{}) (_ js.Value, e error) {
	query = MaybeConvertValueToJs(query)
	return WaitRequest(s.val.Call("getAllKeys", query))
}

// Count counts keys matching the optional query.
func (s *ObjectStore) Count(query interface{}) (_ int, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()
	query = MaybeConvertValueToJs(query)
	v, err := WaitRequest(s.val.Call("count", query))
	if err != nil {
		return 0, err
	}
	return v.Int(), nil
}

// OpenCursor opens a cursor with a optional IDBKeyRange.
// Use Bound() to build a key range.
func (s *ObjectStore) OpenCursor(krv js.Value) (c *Cursor, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()

	req := s.val.Call("openCursor", krv)
	return NewCursor(req), nil
}
