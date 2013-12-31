package irc

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUser(t *testing.T) {
	tests := map[string]*User{
		"foo!bar@irc.localhost": {
			Nickname: "foo",
			Username: "bar",
			Hostname: "irc.localhost",
		},
		"irc.localhost": {
			Nickname: "irc.localhost",
		},
	}

	Convey("NewUser should", t, func() {
		for k, v := range tests {
			Convey(fmt.Sprintf("parse '%s' user", k), func() {
				So(NewUser(k), ShouldResemble, v)
			})
		}
		Convey("compare two users", func() {
			a, b := NewUser("foo!bar@irc.localhost"),
				NewUser("foo!bar@irc.localhost")
			So(a.Equals(b), ShouldBeTrue)
		})
	})
}
