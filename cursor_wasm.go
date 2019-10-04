// +build js,wasm

package indexeddb

import (
	"syscall/js"
)

// Cursor is a object store cursor.
type Cursor struct {
	val        js.Value
	lastCursor js.Value
	nextCh     chan *CursorValue
}

// CursorValue is a object store cursor value.
type CursorValue struct {
	Key   js.Value
	Value js.Value
}

// NewCursor builds a new cursor and registers the onsuccess handler.
func NewCursor(val js.Value) *Cursor {
	c := &Cursor{val: val}
	c.nextCh = make(chan *CursorValue, 1)
	val.Set("onsuccess", js.FuncOf(
		func(th js.Value, dats []js.Value) interface{} {
			cursor := dats[0].Get("target").Get("result")
			global := js.Global()
			c.lastCursor = cursor
			if !cursor.Truthy() {
				close(c.nextCh)
			} else {
				global.Set("resultDebugCursor", cursor)
				cv := &CursorValue{
					Key:   global.Get("Uint8Array").New(cursor.Get("key")),
					Value: cursor.Get("value"),
				}
				go func() {
					c.nextCh <- cv
				}()
			}
			return nil
		},
	))
	return c
}

// WaitValue waits for a value or for the cursor to finish.
// If the cursor is completed, returns nil.
func (c *Cursor) WaitValue() *CursorValue {
	v, ok := <-c.nextCh
	if !ok {
		return nil
	}
	return v
}

// ContinueCursor should be called after WaitValue to trigger a new value to be fetched.
func (c *Cursor) ContinueCursor() {
	c.lastCursor.Call("continue")
}
