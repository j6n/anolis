package irc

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/j6n/logger"
)

var log logger.Logger

func initLogger(b ...bool) {
	if len(b) > 0 && !b[0] {
		log = logger.NewConsoleLogger(logger.INFO)
		return
	}

	// debug logger, really verbose
	// I really, really don't like this
	lines := func() string {
		// 3 because caller -> lambda -> print
		pc, file, line, _ := runtime.Caller(3)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}

		name := runtime.FuncForPC(pc).Name()
		name = name[strings.LastIndex(name, ".")+1:]
		return fmt.Sprintf("%s:%d#%s(): ", short, line, name)
	}

	log = logger.NewConsoleLogger(logger.DEBUG)
	log.SetPrefix(lines)
}
