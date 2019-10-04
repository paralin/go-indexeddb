// +build js,!wasm

package indexeddb

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// WaitRequest waits for an IDBRequest.
// Registers onsuccess and onerror
func WaitRequest(obj *js.Object) (*js.Object, error) {
	ret := func() (*js.Object, error) {
		var err error
		if o := obj.Get("error"); o != nil && o != js.Undefined {
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
	obj.Set("onsuccess", func(e *js.Object) {
		rerr()
	})
	obj.Set("onerror", func(e *js.Object) {
		rerr()
	})
	js.Global.Set("waitTransaction", obj)
	<-errCh
	return ret()
}
