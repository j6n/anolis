package irc

// Conn is an abstract connection
type Conn interface {
	Close()
	WaitForClose() <-chan struct{}
}
