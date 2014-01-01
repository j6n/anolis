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

func (m *MockConn) Join(room string)            { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Part(room string)            { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Kick(room, user, msg string) { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Nick(nick string)            { m.ev.Dispatch(m.msg, m) }
func (m *MockConn) Quit(msg string)             { m.ev.Dispatch(m.msg, m) }

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

func TestConnection(t *testing.T) {
	mock := NewMockConn()
	Convey("connection should", t, func() {
		mock.Do(func() { mock.Join("#hello") }, mock.local, "JOIN", "#hello")
		Convey("add a channel when we join", func() {
			ch, ok := mock.Channels().Get("#hello")
			So(ok, ShouldBeTrue)
			So(ch.Users().Has(mock.local), ShouldBeTrue)
		})
		Convey("remove a channel when we part", func() {
			mock.Do(func() { mock.Part("#hello") }, mock.local, "PART", "#hello", ":byt")
			_, ok := mock.Channels().Get("#hello")
			So(ok, ShouldBeFalse)
		})
	})
}
