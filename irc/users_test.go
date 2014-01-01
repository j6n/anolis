package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUsers(t *testing.T) {
	users := &Users{m: make(map[string]*User)}
	user := NewUser("foo!bar@irc.localhost")

	Convey("users should", t, func() {
		Convey("add a user", func() {
			users.Add(user)
			Convey("and check by name", func() {
				So(users.HasName(user.Nickname), ShouldBeTrue)
			})

			Convey("and check by reference", func() {
				So(users.Has(user), ShouldBeTrue)
			})
		})

		Convey("remove a user", func() {
			Convey("by name", func() {
				users.Add(user)
				users.RemoveName(user.Nickname)
				So(users.HasName(user.Nickname), ShouldBeFalse)
			})

			Convey("by reference", func() {
				users.Add(user)
				users.Remove(user)
				So(users.Has(user), ShouldBeFalse)
			})
		})

		Convey("update a user", func() {
			users.Add(user)
			users.Update("bar", user)
			So(users.HasName("bar"), ShouldBeTrue)
			So(users.HasName("foo"), ShouldBeFalse)
		})
	})
}
