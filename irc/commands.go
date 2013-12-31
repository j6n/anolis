package irc

// Commands represents some common IRC commands
type Commands interface {
	Join(room string)
	Part(room string)
	Kick(room, user, msg string)
	Nick(nick string)
	Quit(msg string)

	Raw(f string, args ...interface{})
	Privmsg(t, f string, args ...interface{})
	Notice(t, f string, args ...interface{})

	Context() Conn
}
