// +build js

package indexeddb

import (
	"errors"
	"syscall/js"
)

// WaitRequest waits for an IDBRequest.
// Registers onsuccess and onerror
func WaitRequest(obj js.Value) (js.Value, error) {
	ret := func() (js.Value, error) {
		var err error
		if o := obj.Get("error"); o.Truthy() {
			err = errors.New(o.Get("message").String())
		}
		return obj.Get("result"), err
	}
	if obj.Get("readyState").String() == "done" {
		return ret()
	}
	errCh := make(chan struct{}, 1)
	rerr := func() {
		select {
		case errCh <- struct{}{}:
		default:
		}
	}
	cb := js.FuncOf(func(th js.Value, dats []js.Value) interface{} {
		go rerr()
		return nil
	})
	obj.Set("onerror", cb)
	obj.Set("onsuccess", cb)
	<-errCh
	return ret()
}
