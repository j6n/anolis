package irc

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockConn struct {
	ch  *Channels
	msg *Message
	ev  *Events

	local, user *User
}

func NewMockConn() *MockConn {
	return &MockConn{
		ch:    &Channels{m: make(map[string]*Channel)},
		local: NewUser("anolis!bot@i.am.a.bot"),
		user:  NewUser("foo!bar@irc.localhost"),
		ev:    NewEvents(),
	}
}

// No-op
func (m *MockConn) Close() {}

// No-op
func (m *MockConn) WaitForClose() <-chan struct{} { return nil }

func (m *MockConn) CurrentNick() string { return m.local.Nickname }
func (m *MockConn) UpdateNick(s string) { m.local.Nickname = s }

func (m *MockConn) Join(room string)    { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Part(room string)    { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Kick(r, u, a string) { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Nick(nick string)    { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Quit(msg string)     { m.ev.Dispatch(m.msg, m) }

func (m *MockConn) Raw(f string, args ...interface{})        { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Privmsg(t, f string, args ...interface{}) { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Notice(t, f string, args ...interface{})  { m.ev.Dispatch(m.msg, m) }

func (m *MockConn) Channels() *Channels { return m.ch }
func (m *MockConn) Connection() Conn    { return m }
func (m *MockConn) Commands() Commands  { return m }

func (m *MockConn) Do(fn func(), u *User, ev string, args ...string) {
	m.msg = ParseMessage(fmt.Sprintf(
		":%s!%s@%s %s %s",
		u.Nickname, u.Username, u.Hostname,
		ev, strings.Join(args, " "),
	))
	m.msg.Source = u
	fn()
}

func TestConnection_LocalUser(t *testing.T) {
	mock := NewMockConn()
	Convey("connection should", t, func() {
		mock.Do(func() { mock.Join("#hello") }, mock.local, "JOIN", "#hello")

		Convey("add a channel when we join", func() {
			ch, ok := mock.Channels().Get("#hello")
			So(ok, ShouldBeTrue)
			So(ch.Users().Has(mock.local), ShouldBeTrue)
		})

		Convey("remove a channel", func() {
			Convey("when we part", func() {
				mock.Do(func() { mock.Part("#hello") }, mock.local, "PART", "#hello", ":byt")
				_, ok := mock.Channels().Get("#hello")
				So(ok, ShouldBeFalse)
			})

			Convey("when we get kicked", func() {
				mock.Do(func() { mock.Kick("#hello", mock.local.Nickname, "bye") },
					mock.user, "KICK", "#hello", mock.local.Nickname, ":bye")
				_, ok := mock.Channels().Get("#hello")
				So(ok, ShouldBeFalse)
			})
		})

		Convey("update our nickname", func() {
			mock.Do(func() { mock.Nick("anolis_") }, mock.local, "NICK", "anolis_")
			So(mock.CurrentNick(), ShouldEqual, "anolis_")
			ch, ok := mock.Channels().Get("#hello")
			So(ok, ShouldBeTrue)
			So(ch.Users().Has(mock.local), ShouldBeTrue)
		})
	})
}

func TestConnection_User(t *testing.T) {
	mock := NewMockConn()
	Convey("connection should update channel", t, func() {
		mock.Do(func() { mock.Join("#hello") }, mock.local, "JOIN", "#hello")
		mock.Do(func() { mock.Join("#hello") }, mock.user, "JOIN", ":#hello")

		Convey("when a user joins", func() {
			ch, ok := mock.Channels().Get("#hello")
			So(ok, ShouldBeTrue)
			So(ch.Users().Has(mock.user), ShouldBeTrue)
		})

		Convey("when a user parts", func() {
			ch, _ := mock.Channels().Get("#hello")
			mock.Do(func() { mock.Part("#hello") }, mock.user, "PART", "#hello", ":bye")
			So(ch.Users().Has(mock.user), ShouldBeFalse)
		})

		Convey("when a user gets kicked", func() {
			ch, _ := mock.Channels().Get("#hello")
			mock.Do(func() { mock.Kick("#hello", mock.user.Nickname, "bye") },
				mock.local, "KICK", "#hello", mock.user.Nickname, ":bye")
			So(ch.Users().Has(mock.user), ShouldBeFalse)
		})

		Convey("when a user quits", func() {
			ch, _ := mock.Channels().Get("#hello")
			mock.Do(func() { mock.Quit("bye") }, mock.user, "QUIT", ":bye")
			So(ch.Users().Has(mock.user), ShouldBeFalse)
		})

		Convey("when a user changes names", func() {
			ch, _ := mock.Channels().Get("#hello")
			mock.Do(func() { mock.Nick("baz") }, mock.user, "NICK", "baz")
			So(mock.user.Nickname, ShouldEqual, "baz")
			So(ch.Users().Has(mock.user), ShouldBeTrue)
		})

		Convey("when the topic changes", func() {
			ch, _ := mock.Channels().Get("#hello")
			mock.Do(func() {
				mock.Raw(":%s!%s@%s TOPIC #hello :test this",
					mock.local.Nickname, mock.local.Username, mock.local.Hostname)
			}, mock.local, "TOPIC", "#hello", ":test this")
			So(ch.GetTopic(), ShouldEqual, "test this")
		})
	})
}

func TestConnection_Channels(t *testing.T) {
	mock := NewMockConn()
	Convey("connection should update channels", t, func() {
		mock.Do(func() { mock.Join("#hello") }, mock.local, "JOIN", "#hello")
		mock.Do(func() { mock.Join("#test") }, mock.local, "JOIN", "#test")
		mock.Do(func() { mock.Join("#world") }, mock.local, "JOIN", "#world")

		mock.Do(func() { mock.Join("#hello") }, mock.user, "JOIN", ":#hello")
		mock.Do(func() { mock.Join("#test") }, mock.user, "JOIN", ":#test")
		mock.Do(func() { mock.Join("#world") }, mock.user, "JOIN", ":#world")

		Convey("when a user quits", func() {
			mock.Do(func() { mock.Quit("bye") }, mock.user, "QUIT", ":bye")
			a, _ := mock.Channels().Get("#hello")
			c, _ := mock.Channels().Get("#test")
			b, _ := mock.Channels().Get("#world")

			So(a.Users().Has(mock.user), ShouldBeFalse)
			So(b.Users().Has(mock.user), ShouldBeFalse)
			So(c.Users().Has(mock.user), ShouldBeFalse)
		})

		Convey("when a user changes nick", func() {
			mock.Do(func() { mock.Nick("baz") }, mock.user, "NICK", "baz")
			a, _ := mock.Channels().Get("#hello")
			c, _ := mock.Channels().Get("#test")
			b, _ := mock.Channels().Get("#world")

			So(mock.user.Nickname, ShouldEqual, "baz")
			So(a.Users().Has(mock.user), ShouldBeTrue)
			So(b.Users().Has(mock.user), ShouldBeTrue)
			So(c.Users().Has(mock.user), ShouldBeTrue)
		})
	})
}
