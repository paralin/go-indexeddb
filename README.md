# IndexedDB

> Basic IndexedDB bindings for GopherJS.

## Getting Started

Check out the example:

```go
import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/paralin/go-indexeddb"
)

db, err := indexeddb.GlobalIndexedDB().Open(ctx, "test-db")
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

Reference:
https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API/Using_IndexedDB

## License

MIT
