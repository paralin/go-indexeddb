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
	// increment ver when you want to change the db structure
	ver := 3
	db, err := indexeddb.GlobalIndexedDB().Open(
		ctx,
		"test-db",
		ver,
		func(d *indexeddb.DatabaseUpdate, oldVersion, newVersion int) error {
			id := "testObjectStore"
			if !d.ContainsObjectStore(id) {
				if err := d.CreateObjectStore(id, nil); err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		fmt.Println("error opening database: " + err.Error())
		return
	}
	js.Global.Set("openedDatabase", db)
	fmt.Println("opened database")
}
