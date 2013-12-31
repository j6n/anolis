package irc

import (
	"log"
	"net/textproto"
	"sync"
)

// Conn is an abstract connection
type Conn interface {
	Close()
	WaitForClose() <-chan struct{}
}

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
func (c *Connection) WaitForClose() <-chan struct{} { return c.done }
