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
goroutines. To fix this, the IndexedDB code in this library takes the extra step
of storing all pending changes in an in-memory write-ahead-log. If the
transaction "goes inactive," the code will re-start the transaction and re-play
all previous actions from the log.

Reference:
https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API/Using_IndexedDB

## License

MIT
