package irc

import (
	"strings"
	"sync"
)

// Event represents an IRC event
type Event func(msg *Message, conn Conn)

// Events holds a collection of events
type Events struct {
	sync.RWMutex
	m map[string][]Event
}

// NewEvents returns a new Events collection
func NewEvents() *Events {
	ev := &Events{m: make(map[string][]Event)}
	return ev
}

// Add register the Event function with the command
func (e *Events) Add(cmd string, fn Event) {
	e.Lock()
	defer e.Unlock()

	c := strings.ToUpper(cmd)
	e.m[c] = append(e.m[c], fn)
}

// Dispatch finds each Event that matches msg's command and invokes it
// with the Message and the Conn
func (e *Events) Dispatch(msg *Message, conn Conn) {
	e.RLock()
	defer e.RUnlock()

	if evs, ok := e.m[strings.ToUpper(msg.Command)]; ok {
		for _, ev := range evs {
			ev(msg, conn)
		}
	}
}
