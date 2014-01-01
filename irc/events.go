package irc

import (
	"log"
	"strings"
	"sync"
)

// Event represents an IRC event
type Event func(msg *Message, ctx Context)

// Events holds a collection of events
type Events struct {
	sync.RWMutex
	m map[string][]Event
}

// NewEvents returns a new Events collection
func NewEvents() *Events {
	ev := &Events{m: make(map[string][]Event)}

	// default events
	ev.Add("JOIN", JoinEvent)
	ev.Add("PART", PartEvent)
	ev.Add("QUIT", QuitEvent)
	ev.Add("KICK", KickEvent)
	ev.Add("NICK", NickEvent)
	ev.Add("TOPIC", TopicEvent)
	ev.Add("ERROR", ErrorEvent)

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
func (e *Events) Dispatch(msg *Message, ctx Context) {
	e.RLock()
	defer e.RUnlock()

	if evs, ok := e.m[strings.ToUpper(msg.Command)]; ok {
		for _, ev := range evs {
			ev(msg, ctx)
		}
	}
}

// Default events

// PingEvent replies to a 'PONG' from the server
func PingEvent(msg *Message, ctx Context) {
	ctx.Commands().Raw("PONG %s", msg.Message)
}

// JoinEvent reacts to the 'JOIN' event
func JoinEvent(msg *Message, ctx Context) {
	if len(msg.Args) > 0 {
		ctx.Channels().Add(msg.Args[0])
		ch, _ := ctx.Channels().Get(msg.Args[0])
		ch.Users().Add(msg.Source)
	} else if ch, ok := ctx.Channels().Get(msg.Message); ok {
		ch.Users().Add(msg.Source)
	}
}

// PartEvent reacts to the 'PART' event
func PartEvent(msg *Message, ctx Context) {}

// KickEvent reacts to the 'KICK' event
func KickEvent(msg *Message, ctx Context) {}

// QuitEvent reacts to the 'QUIT' event
func QuitEvent(msg *Message, ctx Context) {}

// NickEvent reacts to the 'NICK' event
func NickEvent(msg *Message, ctx Context) {}

// TopicEvent reacts to the 'TOPIC' event
func TopicEvent(msg *Message, ctx Context) {}

// ErrorEvent handles any 'ERROR's from the server
func ErrorEvent(msg *Message, ctx Context) {
	log.Fatalln(msg)
}
