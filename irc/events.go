package irc

import (
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
	ev.Add("PING", PingEvent)
	ev.Add("JOIN", JoinEvent)
	ev.Add("PART", PartEvent)
	ev.Add("QUIT", QuitEvent)
	ev.Add("KICK", KickEvent)
	ev.Add("NICK", NickEvent)
	ev.Add("TOPIC", TopicEvent)
	ev.Add("ERROR", ErrorEvent)
	ev.Add("PRIVMSG", PrivmsgEvent)

	return ev
}

// Add register the Event function with the command
func (e *Events) Add(cmd string, fn Event) {
	e.Lock()
	defer e.Unlock()

	c := strings.ToUpper(cmd)
	log.Debugf("adding event: %c", c)
	e.m[c] = append(e.m[c], fn)
}

// Dispatch finds each Event that matches msg's command and invokes it
// with the Message and the Conn
func (e *Events) Dispatch(msg *Message, ctx Context) {
	e.RLock()
	defer e.RUnlock()

	if evs, ok := e.m[strings.ToUpper(msg.Command)]; ok {
		for k, ev := range evs {
			log.Debugf("dispatching to: %s -> %s", k, msg.String())
			ev(msg, ctx)
		}
	}
}

// Default events

// PingEvent replies to a 'PONG' from the server
func PingEvent(msg *Message, ctx Context) {
	ctx.Commands().Raw("PONG %s", msg.Message)
}

// JoinEvent updates the channel user list when a user joins
func JoinEvent(msg *Message, ctx Context) {
	if len(msg.Args) > 0 {
		// it was us that joined
		ctx.Channels().Add(msg.Args[0])
		log.Debugf("creating new channel '%s'", msg.Args[0])
		ch, _ := ctx.Channels().Get(msg.Args[0])
		ch.Users().Add(msg.Source)
		log.Debugf("adding user '%s' to '%s'", msg.Source.Nickname, ch.Name)
	} else if ch, ok := ctx.Channels().Get(msg.Message); ok {
		ch.Users().Add(msg.Source)
		log.Debugf("adding user '%s' to '%s'", msg.Source.Nickname, ch.Name)
	}
}

// PartEvent updates the channel user list when a user leaves
func PartEvent(msg *Message, ctx Context) {
	if ctx.Connection().CurrentNick() == msg.Source.Nickname {
		// it was us that left
		ctx.Channels().Remove(msg.Args[0])
		log.Debugf("removing channel '%s'", msg.Args)
	} else if ch, ok := ctx.Channels().Get(msg.Args[0]); ok {
		ch.Users().Remove(msg.Source)
		log.Debugf("removing '%s' from channel '%s'", msg.Source.Nickname, ch.Name)
	}
}

// KickEvent updates the channel user list when a kick happens
func KickEvent(msg *Message, ctx Context) {
	kickee, room := lastString(msg.Args), msg.Args[0]
	if ctx.Connection().CurrentNick() == kickee {
		// we were kicked
		ctx.Channels().Remove(room)
		log.Debugf("removing channel '%s'", room)
	} else if ch, ok := ctx.Channels().Get(room); ok {
		ch.Users().RemoveName(kickee)
		log.Debugf("removing '%s' from channel '%s'", kickee, room)
	}
}

// QuitEvent updates the channels user list when a user quits
func QuitEvent(msg *Message, ctx Context) {
	ctx.Channels().forEach(msg.Source, func(ch *Channel) {
		ch.Users().Remove(msg.Source)
		log.Debugf("removing '%s' from channel '%s'", msg.Source.Nickname, ch.Name)
	})
}

// NickEvent updates the user list when a nick change happens
func NickEvent(msg *Message, ctx Context) {
	clone, nick := msg.Source.Clone(), msg.Source.Nickname
	msg.Source.Nickname = msg.Args[0]

	// update local, if its us
	if clone.Nickname == ctx.Connection().CurrentNick() {
		ctx.Connection().UpdateNick(msg.Args[0])
	}

	ctx.Channels().forEach(clone, func(ch *Channel) {
		ch.Users().RemoveName(nick)
		ch.Users().Add(msg.Source)
		log.Debugf("changing nick '%s' to '%s' on from channel '%s'", nick, msg.Source.Nickname, ch.Name)
	})
}

// TopicEvent updates the channel topic when it changes
func TopicEvent(msg *Message, ctx Context) {
	if ch, ok := ctx.Channels().Get(msg.Args[0]); ok {
		ch.Topic(msg.Message)
	}
}

// PrivmsgEvent updates the user list when a user speaks
func PrivmsgEvent(msg *Message, ctx Context) {
	if msg.Args[0][0] != 35 {
		return // private message
	}

	if ch, ok := ctx.Channels().Get(msg.Args[0]); ok {
		ch.Users().Update(msg.Source.Nickname, msg.Source)
	}
}

// ErrorEvent handles any 'ERROR's from the server
func ErrorEvent(msg *Message, ctx Context) {
	log.Fatalf(msg.String())
}

// helpers
func lastString(s []string) string { return s[len(s)-1] }
