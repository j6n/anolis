package irc

// Configuration holds the info required to connect
type Configuration struct {
	Hostname string
	Port     int
	Password string

	Nickname, Username, Realname string
	Alternative                  string
}

// NewConfiguration returns a new, default configuration
func NewConfiguration() *Configuration {
	c := &Configuration{
		Hostname: "localhost",
		Port:     6667,
		Nickname: "anolis",
	}

	c.Username = c.Nickname
	c.Realname = c.Nickname
	c.Alternative = c.Nickname + "_"

	return c
}
