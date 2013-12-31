package irc

import (
	"sync"
)

// Channel represents an IRC channel
type Channel struct {
	Name, topic string

	sync.RWMutex
}

// NewChannel creates a new Channel named 'name'
func NewChannel(name string) *Channel {
	return &Channel{Name: name}
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
