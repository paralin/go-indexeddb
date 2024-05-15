# IndexedDB

> Bindings for Go to IndexedDB for WebAssembly and GopherJS.

## Deprecated

It is recommended to use the [go-indexeddb] library instead.

[go-indexeddb]: https://github.com/aperturerobotics/go-indexeddb

## Getting Started

Check out the [example](./example/example.go).

Sample:

```go
  key := []byte("key")
  val := []byte("test")

  err := objStore.Set(key, val)
  err = objStore.Commit()

  data, found, err := objStore.Get(key)
  if err == nil && !found {
    err = errors.New("key not found after setting it")
  }
  // data contains same data as "val"

  prefix := []byte("ke")
  err = objStore.ScanPrefix(prefix, func(key, val []byte) error {
    fmt.Printf("got key/value pair: %v => %v\n", key, val)
    return nil
  })
```

## Transactions expiring

In IndexedDB, transactions will expire if inactive for a short period of time,
or if the Go code goes inactive (such as when waiting for a select statement).
After the transaction expires, all requests will panic / return an error -
"transaction is not active."

This unfortunately happens quite frequently with the Go implementation of the
IndexedDB client in this library, because the Go wasm and/or GopherJS
implementations frequently unwind the stack to the event loop when switching
goroutines. The code will also sometimes panic if the Js code throws any errors.

The IndexedDB code in this library tries to be as minimal of a wrapper around
the underlying JavaScript implementations as possible. As such, the fix for
these issues is implemented in an additional wrapper. After constructing a
`Database`, call NewDurableTransaction(db, scope, mode) instead of Transaction.
If the transaction "goes inactive," it will will re-start the transaction. It
will also handle any panics from the calls.

Unfortunately, a transaction "going inactive" will also commit the transaction.
The "abort" call will "roll-back" the changes made by the transaction. This is a
fairly weak transaction mechanism and should not be relied upon like a
traditional transaction system (in BoltDB or similar).

The "Kvtx" implementation has a easy to use get/set API using `[]byte` slices.
It also implements "ScanPrefix" and "ScanPrefixKeys" for iterating over the db.

Reference:
https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API/Using_IndexedDB

## License

MIT
