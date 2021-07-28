// +build !js

package indexeddb

import (
	"errors"
	"syscall/js"
)

// DurableTransaction handles call panics, errors, and transactions going inactive.
//
// backs any writes with an in-memory write-ahead-log
type DurableTransaction struct {
	// d is the database
	d *Database
	// txn is the current transaction
	// note: if transaction is inactive, this becomes nil
	txn *Transaction
	// scope is the scope
	scope []string
	// mode is the txn mode
	mode TransactionMode
	// stores is the set of object store handles
	stores map[string]*DurableObjectStore
}

// NewDurableTransaction starts a transaction that handles typical errors and panics.
//
// This is the recommended way to use this library.
func NewDurableTransaction(d *Database, scope []string, mode TransactionMode) (*DurableTransaction, error) {
	txn, err := d.Transaction(scope, mode)
	if err != nil {
		return nil, err
	}
	mode = txn.GetMode()

	dt := &DurableTransaction{
		d:      d,
		txn:    txn,
		scope:  scope,
		mode:   mode,
		stores: make(map[string]*DurableObjectStore),
	}
	dt.setOnCompleteCallback()
	return dt, nil
}

// setOnCompleteCallback sets the on-complete callback.
func (t *DurableTransaction) setOnCompleteCallback() {
	if t.txn != nil {
		t.txn.val.Set(
			"oncomplete",
			js.FuncOf(func(th js.Value, dats []js.Value) interface{} {
				// set txn to nil to indicate transaction complete
				t.txn = nil
				return nil
			}),
		)
	}
}

// restartTransaction restarts the tx
func (t *DurableTransaction) restartTransaction() error {
	txn, err := t.d.Transaction(t.scope, t.mode)
	if err != nil {
		return err
	}
	for sid, stor := range t.stores {
		nstor, err := txn.GetObjectStore(sid)
		if err != nil {
			return err
		}
		stor.store = nstor
		ops := stor.ops
		for i, op := range ops {
			if err := op.apply(nstor); err != nil {
				return err
			}
			stor.ops = ops[i+1:] // don't apply again if successful
		}
		if len(stor.ops) == 0 {
			stor.ops = nil
		}
	}
	t.txn = txn
	t.setOnCompleteCallback()
	return nil
}

// GetMode returns the transaction mode.
func (t *DurableTransaction) GetMode() TransactionMode {
	return t.mode
}

// Abort aborts a transaction.
func (t *DurableTransaction) Abort() {
	if t.txn != nil {
		t.txn.Abort()
		t.txn = nil
	}
}

// Commit commits a transaction and waits for it to complete
func (t *DurableTransaction) Commit() error {
	var err error
	attempts := 0
	// force restart txn for commit
	for t.txn == nil {
		attempts++
		if attempts > 10 {
			if err == nil {
				err = errors.New("unable to restart transaction without it going inactive")
			}
			return err
		}
		// restartTransaction also flushes all pending ops.
		if err = t.restartTransaction(); err != nil {
			if errIsInactiveTransaction(err) {
				continue
			} else {
				return err
			}
		}
	}
	if t.txn != nil {
		t.txn.Commit()
		err = t.txn.WaitComplete()
		t.txn = nil
	}
	return err
}

// Restart restarts the transaction if inactive.

// DurableObjectStore backs changes in a write-ahead log.
type DurableObjectStore struct {
	id string
	tx *DurableTransaction
	// ops is the operation log so far
	ops []*durableOp
	// store may become nil if the transaction is inactive
	store *ObjectStore
}

// GetObjectStore returns a object store.
func (t *DurableTransaction) GetObjectStore(id string) (o *DurableObjectStore, e error) {
	if s, ok := t.stores[id]; ok {
		return s, nil
	}
	var store *ObjectStore
	if t.txn != nil {
		var err error
		store, err = t.txn.GetObjectStore(id)
		if err != nil {
			if errIsInactiveTransaction(err) {
				t.txn = nil
				store = nil
			} else {
				return nil, err
			}
		}
	}
	o = &DurableObjectStore{
		id:    id,
		tx:    t,
		store: store,
	}
	t.stores[id] = o
	return o, nil
}

// GetName returns the object store name.
func (s *DurableObjectStore) GetName() string {
	return s.id
}

// durableOp contains an operation against the store.
type durableOp struct {
	// apply applies the op
	apply func(s *ObjectStore) error
}

// newDurableOp constructs a new durableOp
func newDurableOp(apply func(s *ObjectStore) error) *durableOp {
	return &durableOp{apply: apply}
}

// pushOp attempts an operation with the "inactive transaction" logic
func (s *DurableObjectStore) pushOp(op *durableOp) error {
	if s.tx.txn != nil && s.store != nil {
		err := op.apply(s.store)
		if err != nil && errIsInactiveTransaction(err) {
			s.tx.txn = nil
			s.store = nil
			err = nil
			// defer applying the op until Commit() or a read
			s.ops = append(s.ops, op)
		}
		if err != nil {
			return err
		}
	} else {
		s.ops = append(s.ops, op)
	}
	return nil
}

// getOrBuildStore gets the store or re-starts the tx if it's inactive
func (s *DurableObjectStore) getOrBuildStore() (*ObjectStore, error) {
	if s.tx.txn != nil && s.store != nil {
		return s.store, nil
	}
	if err := s.tx.restartTransaction(); err != nil {
		return nil, err
	}
	ns := s.store
	if ns == nil {
		return nil, errors.New("indexed-db store is nil after restarting tx")
	}
	return ns, nil
}

// Put puts data into the store.
func (s *DurableObjectStore) Put(value interface{}, key interface{}) error {
	value = MaybeConvertValueToJs(value)
	key = MaybeConvertValueToJs(key)
	return s.pushOp(newDurableOp(func(s *ObjectStore) error {
		return s.Put(value, key)
	}))
}

// Add adds data to the store.
func (s *DurableObjectStore) Add(value interface{}, key interface{}) error {
	value = MaybeConvertValueToJs(value)
	key = MaybeConvertValueToJs(key)
	return s.pushOp(newDurableOp(func(s *ObjectStore) error {
		return s.Add(value, key)
	}))
}

// Delete deletes data from the store.
func (s *DurableObjectStore) Delete(query interface{}) error {
	query = MaybeConvertValueToJs(query)
	return s.pushOp(newDurableOp(func(s *ObjectStore) error {
		return s.Delete(query)
	}))
}

// Clear clears all data from the store.
func (s *DurableObjectStore) Clear() error {
	return s.pushOp(newDurableOp(func(s *ObjectStore) error {
		return s.Clear()
	}))
}

// durableRead retries a read several times upon "inactive transaction" errors
func (s *DurableObjectStore) durableRead(read func(stor *ObjectStore) (js.Value, error)) (js.Value, error) {
	attempts := 0
	for {
		stor, err := s.getOrBuildStore()
		if err != nil {
			return js.Undefined(), err
		}
		result, err := read(stor)
		if err != nil {
			attempts++
			if errIsInactiveTransaction(err) {
				s.tx.txn = nil
				s.store = nil
				// retry
				if attempts > 10 {
					return js.Undefined(), err
				} else {
					continue
				}
			}
		}
		return result, nil
	}
}

// Get gets data from the store
func (s *DurableObjectStore) Get(query interface{}) (js.Value, error) {
	return s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		return stor.Get(query)
	})
}

// GetKey gets data from the store by key.
func (s *DurableObjectStore) GetKey(query interface{}) (js.Value, error) {
	return s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		return stor.GetKey(query)
	})
}

// GetAll gets all values matching an optional query with an optional count.
func (s *DurableObjectStore) GetAll(query interface{}) (js.Value, error) {
	return s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		return stor.GetAll(query)
	})
}

// GetAllKeys gets all keys matching an optional query with an optional count.
func (s *DurableObjectStore) GetAllKeys(query interface{}) (js.Value, error) {
	return s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		return stor.GetAllKeys(query)
	})
}

// Count counts keys matching the optional query.
func (s *DurableObjectStore) Count(query interface{}) (int, error) {
	var out int
	_, err := s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		return stor.GetAllKeys(query)
		c, err := stor.Count(query)
		out = c
		return js.Undefined(), err
	})
	return out, err
}

// OpenCursor opens a cursor with a optional IDBKeyRange.
// Use Bound() to build a key range.
func (s *DurableObjectStore) OpenCursor(krv js.Value) (c *Cursor, e error) {
	defer func() {
		if err := recover(); err != nil {
			e, _ = err.(error)
		}
	}()

	var out *Cursor
	_, err := s.durableRead(func(stor *ObjectStore) (js.Value, error) {
		c, err := stor.OpenCursor(krv)
		out = c
		return js.Undefined(), err
	})
	return out, err
}
