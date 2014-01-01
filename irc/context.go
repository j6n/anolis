package irc

// Context is the abstract context
type Context interface {
	Channels() *Channels
	Connection() Conn
	Commands() Commands
}
