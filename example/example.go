//+build js

package main

import (
	"context"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/paralin/go-indexeddb"
)

func main() {
	fmt.Println("opening db")
	ctx := context.Background()
	db, err := indexeddb.GlobalIndexedDB().Open(ctx, "test-db")
	if err != nil {
		panic(err)
	}
	js.Global.Set("openedDatabase", db)
	fmt.Println("opened database")
}
