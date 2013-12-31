package irc

import (
	"log"
	"strings"
	"sync"
)

// Event represents an IRC event
type Event func(msg *Message, cmd Commands)

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
func (e *Events) Dispatch(msg *Message, cmd Commands) {
	e.RLock()
	defer e.RUnlock()

	if evs, ok := e.m[strings.ToUpper(msg.Command)]; ok {
		for _, ev := range evs {
			ev(msg, cmd)
		}
	}
}

// Default events

// PingEvent replies to a 'PONG' from the server
func PingEvent(msg *Message, cmd Commands) {
	cmd.Raw("PONG %s", msg.Message)
}

// JoinEvent reacts to the 'JOIN' event
func JoinEvent(msg *Message, cmd Commands) {}

// PartEvent reacts to the 'PART' event
func PartEvent(msg *Message, cmd Commands) {}

// KickEvent reacts to the 'KICK' event
func KickEvent(msg *Message, cmd Commands) {}

// QuitEvent reacts to the 'QUIT' event
func QuitEvent(msg *Message, cmd Commands) {}

// NickEvent reacts to the 'NICK' event
func NickEvent(msg *Message, cmd Commands) {}

// TopicEvent reacts to the 'TOPIC' event
func TopicEvent(msg *Message, cmd Commands) {}

// ErrorEvent handles any 'ERROR's from the server
func ErrorEvent(msg *Message, cmd Commands) {
	log.Fatalln(msg)
}
