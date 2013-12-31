package irc

import (
	"strings"
)

// ParseMessage takes a raw line and returns a new Message
func ParseMessage(raw string) *Message {
	msg := &Message{Raw: raw}
	// :source command [args] :message
	if raw[0] == ':' {
		if i := strings.Index(raw, " "); i >= -1 {
			msg.Source = NewUser(raw[1:i])
			raw = raw[i+1 : len(raw)]
		}
	}

	args := strings.SplitN(raw, " :", 2)
	if len(args) > 1 {
		msg.Message = args[1]
	}

	args = strings.Split(args[0], " ")
	msg.Command, msg.Args = strings.ToUpper(args[0]), []string{}

	if len(args) > 1 {
		msg.Args = args[1:len(args)]
	}

	return msg
}
