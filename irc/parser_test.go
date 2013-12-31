package irc

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseMessage(t *testing.T) {
	src := NewUser("foo!bar@irc.localhost")
	tests := map[string]*Message{
		"local join": {
			Raw:     ":foo!bar@irc.localhost JOIN #foobar",
			Source:  src,
			Command: "JOIN",
			Args:    []string{"#foobar"},
		},
		"join": {
			Raw:     ":foo!bar@irc.localhost JOIN :#foobar",
			Source:  src,
			Command: "JOIN",
			Args:    []string{"#foobar"},
		},
		"privmsg": {
			Raw:     ":foo!bar@irc.localhost PRIVMSG #foobar :hello world",
			Source:  src,
			Command: "PRIVMSG",
			Args:    []string{"#foobar"},
			Message: "hello world",
		},
		"part": {
			Raw:     ":foo!bar@irc.localhost PART #foobar :bye",
			Source:  src,
			Command: "PART",
			Args:    []string{"#foobar"},
			Message: "bye",
		},
		"raw": {
			Raw:     "NOTICE AUTH :*** Checking Ident",
			Command: "NOTICE",
			Args:    []string{"AUTH"},
			Message: "*** Checking Ident",
		},
		"ping": {
			Raw:     "PING :1594198849",
			Command: "PING",
			Args:    []string{},
			Message: "1594198849",
		},
		"no text": {
			Raw:     ":irc.localhost 004 museun irc.localhost beware1.6.2 dgikoswx biklmnoprstv",
			Source:  NewUser("irc.localhost"),
			Command: "004",
			Args:    []string{"museun", "irc.localhost", "beware1.6.2", "dgikoswx", "biklmnoprstv"},
			Message: "",
		},
		"many colons": {
			Raw:     ":foo!bar@irc.localhost PRIVMSG hello :hello world :) for more :colons",
			Source:  src,
			Command: "PRIVMSG",
			Args:    []string{"hello"},
			Message: "hello world :) for more :colons",
		},
	}

	Convey("Parse Message should", t, func() {
		for k, v := range tests {
			Convey(fmt.Sprintf("parse '%s' message", k), func() {
				So(ParseMessage(v.Raw), ShouldResemble, v)
			})
		}
	})
}
