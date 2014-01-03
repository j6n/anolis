package main

import (
	"log"

	"github.com/j6n/anolis/irc"
)

func main() {
	conf := irc.NewConfiguration()
	// turn on verbose logging
	conf.Verbose = true
	conn := irc.Dial(conf)

	conn.Events().Add("001", func(m *irc.Message, ctx irc.Context) {
		ctx.Commands().Join("#test")
	})

	conn.Events().Add("PRIVMSG", func(m *irc.Message, ctx irc.Context) {
		log.Println(m)
	})

	<-conn.Connection().WaitForClose()
}
