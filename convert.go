// +build js,!wasm

package indexeddb

import "github.com/gopherjs/gopherjs/js"

// Bound builds a new IDBKeyRange with the range.
func Bound(lower, upper interface{}, lowerOpen, upperOpen bool) *js.Object {
	return js.Global.
		Get("IDBKeyRange").
		Call(
			"bound",
			lower,
			upper,
			lowerOpen,
			upperOpen,
		)
}

// MaybeConvertValueToJs conditionally converts val to javascript.
func MaybeConvertValueToJs(val interface{}) interface{} {
	return val
}

// CopyByteSliceToJS copies a byte slice to javascript.
func CopyByteSliceToJs(vb []byte) *js.Object {
	return js.MakeWrapper(vb)
}

// CopyByteSliceFromJS copies a byte slice from javascript.
func CopyByteSliceFromJs(vb *js.Object) []byte {
	return vb.Interface().([]byte)
}
