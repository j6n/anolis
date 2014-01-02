package irc

import "sync"

// Channels represents a collection of channels
type Channels struct {
	m map[string]*Channel
	sync.RWMutex
}

// Add adds a new channel, with the name 'name' to the collection
func (c *Channels) Add(name string) {
	c.Lock()
	defer c.Unlock()

	log.Debugf("adding channel '%s'", name)
	c.m[name] = NewChannel(name) // TODO cache this
}

// Remove removes channel, with the name 'name' from the collection
func (c *Channels) Remove(name string) {
	c.Lock()
	defer c.Unlock()

	log.Debugf("remove channel '%s'", name)
	delete(c.m, name)
}

// Has returns whether the collection has a channel named 'name'
func (c *Channels) Has(name string) (ok bool) {
	_, ok = c.Get(name)
	return
}

// Get returns a channel, named 'name' and if it exists
func (c *Channels) Get(name string) (ch *Channel, ok bool) {
	c.RLock()
	defer c.RUnlock()

	ch, ok = c.m[name]
	return
}

// GetNames returns all of the names of the channels
func (c *Channels) GetNames() []string {
	c.RLock()
	defer c.RUnlock()

	out := make([]string, 0, len(c.m))
	for k := range c.m {
		out = append(out, k)
	}
	return out
}

func (c *Channels) forEach(user *User, fn func(ch *Channel)) {
	for _, name := range c.GetNames() {
		if ch, _ := c.Get(name); ch.Users().Has(user) {
			fn(ch)
		}
	}
}
