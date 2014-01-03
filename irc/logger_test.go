package irc

import (
	"testing"

	"github.com/j6n/logger"
)

type NoopLogger struct{}

func (n *NoopLogger) Debugf(a string, v ...interface{}) {}
func (n *NoopLogger) Infof(a string, v ...interface{})  {}
func (n *NoopLogger) Warnf(a string, v ...interface{})  {}
func (n *NoopLogger) Fatalf(a string, v ...interface{}) {}
func (n *NoopLogger) Errorf(a string, v ...interface{}) error {
	return nil
}

func (n *NoopLogger) SetVerbosity(lvl logger.LogLevel) {}
func (n *NoopLogger) Verbosity() logger.LogLevel {
	return logger.NONE
}

func (n *NoopLogger) SetPrefix(func() string) {}

func TestLogger(t *testing.T) {
	// no tests yet
	initLogger(true)

	world := struct {
		name string
		some int
	}{"world", 1004}

	log.Debugf("hello %+v", world)
	log.Infof("some info %X", 3735928559)

	func() { log.Warnf("a warning :(") }()
}
