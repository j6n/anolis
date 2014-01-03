package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChannels(t *testing.T) {
	log = &NoopLogger{}

	ch := &Channels{m: make(map[string]*Channel)}
	Convey("channels should", t, func() {
		Convey("add a channel", func() {
			ch.Add("#foo")
			So(ch.Has("#foo"), ShouldBeTrue)
		})

		Convey("remove a channel", func() {
			ch.Remove("#foo")
			So(ch.Has("#foo"), ShouldBeFalse)
		})

		Convey("list all channels", func() {
			ch.Add("#foo")
			ch.Add("#bar")

			result := ch.GetNames()
			So(result, ShouldContain, "#foo")
			So(result, ShouldContain, "#bar")
		})
	})
}
