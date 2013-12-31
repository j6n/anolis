package irc

import "strings"

// User represents an IRC user
type User struct {
	Nickname string
	Username string
	Hostname string
}

// NewUser parses a raw string and returns a user
func NewUser(raw string) *User {
	// TODO cache this
	if strings.Index(raw, "!") == -1 {
		return &User{Nickname: raw}
	}

	first := strings.Split(raw, "!")
	nick, second := first[0], strings.Split(first[1], "@")
	user, host := second[0], second[1]

	return &User{nick, user, host}
}

func (u *User) String() string {
	return u.Nickname
}

// Clone returns a copy of the User
func (u *User) Clone() *User {
	return &User{u.Nickname, u.Username, u.Hostname}
}

// Equals compares the User to another User
func (u *User) Equals(other *User) bool {
	return u.Nickname == other.Nickname
}
