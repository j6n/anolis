package irc

import (
	"strings"

	bs "github.com/j6n/bufferedstring"
)

// Message represents an IRC message
type Message struct {
	Raw     string
	Source  string // optional
	Command string
	Args    []string // 15 max
	Message string   // optional
}

func (m *Message) String() string {
	buf := bs.NewBufferedString()
	buf.Add("[%s]", m.Command)
	if m.Source != nil {
		buf.Add("%s", m.Source)
	}
	if len(m.Args) > 0 {
		buf.Add("(%s)", strings.Join(m.Args, ", "))
	}

	buf.Append(m.Message)
	return buf.String()
}
