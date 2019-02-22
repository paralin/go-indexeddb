package indexeddb

import (
	"github.com/gopherjs/gopherjs/js"
)

// Cursor is a object store cursor.
type Cursor struct {
	*js.Object

	lastCursor *js.Object
	nextCh     chan *CursorValue
}

// CursorValue is a object store cursor value.
type CursorValue struct {
	Key   *js.Object
	Value *js.Object
}

// NewCursor builds a new cursor and registers the onsuccess handler.
func NewCursor(req *js.Object) *Cursor {
	c := &Cursor{Object: req}
	c.nextCh = make(chan *CursorValue, 1)
	c.Set("onsuccess", func(e *js.Object) {
		cursor := e.Get("target").Get("result")
		c.lastCursor = cursor
		if cursor == nil || cursor == js.Undefined {
			close(c.nextCh)
		} else {
			js.Global.Set("resultDebugCursor", cursor)
			cv := &CursorValue{
				Key:   js.Global.Get("Uint8Array").New(cursor.Get("key")),
				Value: cursor.Get("value"),
			}
			go func() {
				c.nextCh <- cv
			}()
		}
	})
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
