package irc

// Context is the abstract context
type Context interface {
	Channels() *Channels
	Events() *Events
	Connection() Conn
	Commands() Commands
}
