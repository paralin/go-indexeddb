package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"

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
	// js.Global().Set("openedDatabase", db)
	fmt.Println("opened database")

	durTx, err := indexeddb.NewDurableTransaction(db, []string{id}, indexeddb.READWRITE)
	if err != nil {
		fmt.Println("error getting durable transaction: " + err.Error())
		return
	}

	objStore, err := indexeddb.NewKvtxTx(durTx, id)
	if err != nil {
		fmt.Println("error getting object store: " + err.Error())
		return
	}

	key := []byte("key")
	val := []byte("test")
	if err := objStore.Set(key, val); err != nil {
		fmt.Println(err.Error())
		return
	}

	_ = objStore.Commit()

	dat, found, err := objStore.Get(key)
	if err == nil && !found {
		err = errors.New("key not found after setting it")
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("read data: %v\n", dat)

	prefix := []byte("ke")
	err = objStore.ScanPrefix(prefix, func(key, val []byte) error {
		fmt.Printf("got key/value pair: %v => %v\n", key, val)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// prevent "go program already has exited" error
	if runtime.GOOS == "js" {
		<-ctx.Done()
	}
}
