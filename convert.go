// +build js

package indexeddb

import "syscall/js"

// Bound builds a new IDBKeyRange with the range.
func Bound(lower, upper interface{}, lowerOpen, upperOpen bool) js.Value {
	return js.Global().
		Get("IDBKeyRange").
		Call(
			"bound",
			MaybeConvertValueToJs(lower),
			MaybeConvertValueToJs(upper),
			lowerOpen,
			upperOpen,
		)
}

// MaybeConvertValueToJs conditionally converts val to javascript.
func MaybeConvertValueToJs(val interface{}) interface{} {
	switch vb := val.(type) {
	case []byte:
		return CopyByteSliceToJs(vb)
	}
	return val
}

// CopyByteSliceToJS copies a byte slice to javascript.
func CopyByteSliceToJs(vb []byte) js.Value {
	vba := js.Global().Get("Uint8Array").New(len(vb))
	js.CopyBytesToJS(vba, vb)
	return vba
}

// CopyByteSliceFromJS copies a byte slice from javascript.
func CopyByteSliceFromJs(vb js.Value) []byte {
	b := make([]byte, vb.Length())
	js.CopyBytesToGo(b, vb)
	return b
}
