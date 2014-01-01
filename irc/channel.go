package irc

import (
	"sync"
)

// Channel represents an IRC channel
type Channel struct {
	Name, topic string
	users       *Users

	sync.RWMutex
}

// NewChannel creates a new Channel named 'name'
func NewChannel(name string) *Channel {
	return &Channel{
		Name:  name,
		users: &Users{m: make(map[string]*User)},
	}
}

// Users returns the user collection for the channel
func (c *Channel) Users() *Users {
	return c.users
}

// Topic sets the channels' topic to 't'
func (c *Channel) Topic(t string) {
	c.Lock()
	defer c.Unlock()

	c.topic = t
}

// GetTopic returns the channels' topic
func (c *Channel) GetTopic() string {
	c.RLock()
	defer c.RUnlock()

	return c.topic
}
