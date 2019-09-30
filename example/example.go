package main

import (
	"context"
	"fmt"

	"github.com/paralin/go-indexeddb"
	"syscall/js"
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
	js.Global().Set("openedDatabase", db.GetJsValue())
	fmt.Println("opened database")

	tx, err := db.Transaction([]string{id}, indexeddb.READWRITE)
	if err != nil {
		fmt.Println("error getting transaction: " + err.Error())
		return
	}

	js.Global().Set("openedTransaction", tx.GetJsValue())
	objStore, err := tx.GetObjectStore(id)
	if err != nil {
		fmt.Println("error getting obj store: " + err.Error())
		return
	}

	key := []byte("key")
	val := []byte("test")
	if err := objStore.Put(val, key); err != nil {
		fmt.Println(err.Error())
		return
	}

	dat, err := objStore.Get(key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("read data: %v\n", dat.String())

	prefix := []byte("ke")
	prefixGreater := make([]byte, len(prefix)+1)
	copy(prefixGreater, prefix)
	prefixGreater[len(prefixGreater)-1] = ^byte(0)
	krv := indexeddb.Bound(prefix, prefixGreater, false, false)
	cursor, err := objStore.OpenCursor(krv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cval := cursor.WaitValue()
	cvalKey := indexeddb.CopyByteSliceFromJs(cval.Key)
	cvalVal := indexeddb.CopyByteSliceFromJs(cval.Value)
	fmt.Printf(
		"got value from cursor: key %#v value %#v\n",
		string(cvalKey),
		string(cvalVal),
	)
	cursor.ContinueCursor()

	cval = cursor.WaitValue()
	if cval != nil {
		fmt.Println("expected cval to be nil but it wasn't")
		return
	}

	dat, err = objStore.Get([]byte("notexist"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
