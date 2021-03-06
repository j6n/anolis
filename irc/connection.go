package irc

import (
	"fmt"
	"net/textproto"
	"sync"
)

// Connection represents a connection to an irc server
type Connection struct {
	address  string
	nickname string

	ev *Events
	ch *Channels

	conn *textproto.Conn
	once sync.Once
	done chan struct{}

	sync.RWMutex // TODO use this to make us thread-safe
}

// Dial connects to the address with the nickname
// and returns a Conn
func Dial(conf *Configuration) Context {
	initLogger(conf.Verbose)
	conn := &Connection{
		address:  fmt.Sprintf("%s:%d", conf.Hostname, conf.Port),
		nickname: conf.Nickname,

		ch: &Channels{m: make(map[string]*Channel)},
		ev: NewEvents(),

		done: make(chan struct{}),
	}

	log.Debugf("connecting to %s", conn.address)
	tp, err := textproto.Dial("tcp", conn.address)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conn.conn = tp
	if conf.Password != "" {
		conn.Raw("PASS %s", conf.Password)
	}

	conn.Raw("NICK %s", conf.Nickname)
	conn.Raw("USER %s 0 * :%s", conf.Username, conf.Realname)

	log.Debugf("starting readLoop")
	go conn.readLoop()
	return conn
}

// Close closes the connection
func (c *Connection) Close() {
	c.once.Do(func() {
		log.Debugf("closing connection")
		c.conn.Close()
		close(c.done)
	})
}

// WaitForClose returns a channel that'll be closed when the connection closes
func (c *Connection) WaitForClose() <-chan struct{} {
	return c.done
}

// CurrentNick returns the local users' nickname
func (c *Connection) CurrentNick() string {
	return c.nickname
}

// UpdateNick updates the connections nickname
func (c *Connection) UpdateNick(s string) {
	log.Debugf("updating nickname to: %s", s)
	c.nickname = s
}

// Join sends the join command for room
func (c *Connection) Join(room string) {
	c.Raw("JOIN %s", room)
}

// Part sends the part command for room
func (c *Connection) Part(room string) {
	c.Raw("PART %s", room)
}

// Kick sends the kick command for user on room with msg
func (c *Connection) Kick(room, user, msg string) {
	c.Raw("KICK %s %s :%s", room, user, msg)
}

// Nick sends the nick command with the new nick
func (c *Connection) Nick(nick string) {
	c.Raw("NICK %s", nick)
}

// Quit sends the quit command with msg
func (c *Connection) Quit(msg string) {
	c.Raw("QUIT :%s", msg)
}

// Privmsg sends a private message, f, formatted with args to t
func (c *Connection) Privmsg(t, f string, args ...interface{}) {
	c.Raw("PRIVMSG %s :%s", fmt.Sprintf(f, args...))
}

// Notice sends a notice message, f, formatted with args to t
func (c *Connection) Notice(t, f string, args ...interface{}) {
	c.Raw("NOTICE %s :%s", fmt.Sprintf(f, args...))
}

// Raw sends a raw message, f, formatted with args
func (c *Connection) Raw(f string, args ...interface{}) {
	log.Debugf(">> "+f, args...)
	c.conn.Cmd(f, args...)
}

// Connection returns the Connection's context (this pointer)
func (c *Connection) Connection() Conn {
	return c
}

// Channels returns the Connection's channels
func (c *Connection) Channels() *Channels {
	return c.ch
}

// Commands returns the Connection's commands
func (c *Connection) Commands() Commands {
	return c
}

// Events returns the Connection's events
func (c *Connection) Events() *Events {
	return c.ev
}

func (c *Connection) readLoop() {
	for {
		line, err := c.conn.ReadLine()
		if err != nil {
			c.Close()
			log.Fatalf(err.Error())
			break
		}

		msg := ParseMessage(line)
		log.Debugf("<< %s", msg)
		go c.ev.Dispatch(msg, c)
	}
}
