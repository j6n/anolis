package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockConn struct {
	msg *Message
}

func NewMockConn() *MockConn {
	return &MockConn{}
}

// No-op
func (m *MockConn) Close() {}

// No-op
func (m *MockConn) WaitForClose() <-chan struct{} { return nil }

func (m *MockConn) Join(room string)            {}
func (m *MockConn) Part(room string)            {}
func (m *MockConn) Kick(room, user, msg string) {}
func (m *MockConn) Nick(nick string)            {}
func (m *MockConn) Quit(msg string)             {}

func (m *MockConn) Raw(f string, args ...interface{})        {}
func (m *MockConn) Privmsg(t, f string, args ...interface{}) {}
func (m *MockConn) Notice(t, f string, args ...interface{})  {}

func (m *MockConn) Do(fn func()) {}

func TestConnection(t *testing.T) {
	//mock := NewMockConn()
	Convey("", t, func() {})
}
