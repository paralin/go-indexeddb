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
	id := "testObjectStore"
	db, err := indexeddb.GlobalIndexedDB().Open(
		ctx,
		"test-db",
		ver,
		func(d *indexeddb.DatabaseUpdate, oldVersion, newVersion int) error {
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

	tx, err := db.Transaction([]string{id}, indexeddb.READWRITE)
	if err != nil {
		fmt.Println("error getting transaction: " + err.Error())
		return
	}

	objStore, err := tx.GetObjectStore(id)
	if err != nil {
		fmt.Println("error getting obj store: " + err.Error())
		return
	}

	if err := objStore.Put("value", "key"); err != nil {
		fmt.Println(err.Error())
		return
	}

	dat, err := objStore.Get("key")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("read data: %#v\n", dat.Interface())
}
