package irc

import "sync"

// Users represents a collection of users
type Users struct {
	m map[string]*User
	sync.RWMutex
}

// Has returns whether user is contained in the collection
func (u *Users) Has(user *User) bool {
	return u.HasName(user.Nickname)
}

// HasName returns whether user with name is in the collection
func (u *Users) HasName(name string) bool {
	u.RLock()
	defer u.RUnlock()

	_, ok := u.m[name]
	return ok
}

// Get returns a user, by name and if they were found
func (u *Users) Get(name string) (user *User, ok bool) {
	u.RLock()
	defer u.RUnlock()

	user, ok = u.m[name]
	return
}

// Add adds the user to the collection
func (u *Users) Add(user *User) {
	u.Lock()
	defer u.Unlock()

	u.m[user.Nickname] = user
}

// Update updates User with the new nick
func (u *Users) Update(nick string, user *User) {
	u.Lock()
	defer u.Unlock()

	u.m[nick] = user // TODO cache this
}

// Remove removes the user from the collection
func (u *Users) Remove(user *User) {
	u.RemoveName(user.Nickname)
}

// RemoveName removes the user by name from the collection
func (u *Users) RemoveName(name string) {
	u.Lock()
	defer u.Unlock()

	delete(u.m, name)
}
