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

## License

MIT
