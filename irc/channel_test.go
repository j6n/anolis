package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChannel(t *testing.T) {
	user := NewUser("foo!bar@irc.localhost")
	ch := NewChannel("#test")
	Convey("channel should", t, func() {
		Convey("add a user", func() {
			ch.Users().Add(user)
			So(ch.Users().HasName("foo"), ShouldBeTrue)
		})
		Convey("remove a user", func() {
			ch.Users().Add(user)
			ch.Users().Remove(user)
			So(ch.Users().Has(user), ShouldBeFalse)
		})
		Convey("remove a user by name", func() {
			ch.Users().Add(user)
			ch.Users().RemoveName(user.Nickname)
			So(ch.Users().HasName("foo"), ShouldBeFalse)
		})
	})
}
