package irc

import (
	"fmt"
	"log"
	"net/textproto"
	"sync"
)

// Connection represents a connection to an irc server
type Connection struct {
	address string
	nick    string

	conn *textproto.Conn
	once sync.Once
	done chan struct{}
}

// Dial connects to the address with the nickname
// and returns a Conn
func Dial(address, nickname string) Conn {
	conn := &Connection{
		address:  address,
		nickname: nickname,
		done:     make(chan struct{}),
	}

	tp, err := textproto.Dial("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	conn.conn = tp
	return conn
}

// Close closes the connection
func (c *Connection) Close() {
	c.once.Do(func() {
		c.conn.Close()
		close(c.done)
	})
}

// WaitForClose returns a channel that'll be closed when the connection closes
func (c *Connection) WaitForClose() <-chan struct{} {
	return c.done
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
	c.Quit("QUIT :%s", msg)
}

// Raw sends a raw message, f, formatted with args
func (c *Connection) Raw(f string, args ...interface{}) {
	c.conn.Cmd(f, args...)
}

// Privmsg sends a private message, f, formatted with args to t
func (c *Connection) Privmsg(t, f string, args ...interface{}) {
	c.Raw("PRIVMSG %s :%s", fmt.Sprintf(f, args...))
}

// Notice sends a notice message, f, formatted with args to t
func (c *Connection) Notice(t, f string, args ...interface{}) {
	c.Raw("NOTICE %s :%s", fmt.Sprintf(f, args...))
}
