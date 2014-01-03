package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/j6n/anolis/irc"
)

func main() {
	// to capture Ctrl-C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// create a new configuration
	conf := irc.NewConfiguration()
	// turn on verbose logging
	// conf.Verbose = true

	// create and connect to irc with our configuration
	conn := irc.Dial(conf)

	// join channels when we are fully connected
	conn.Events().Add("001", func(m *irc.Message, ctx irc.Context) {
		ctx.Commands().Join("#test")
	})

	// log when we receive a private emssage
	conn.Events().Add("PRIVMSG", func(m *irc.Message, ctx irc.Context) {
		log.Println(m)
	})

	go func() {
		<-quit
		// send the quit
		conn.Commands().Quit("bye")
	}()

	// wait for the connection to close
	<-conn.Connection().WaitForClose()
}
