package copter

import (
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

// NewLogger returns a new Logger instance
func Logger() HandlerFunc {
	l := &logger{log.New(os.Stdout, "[copter] ", 0)}

	return func(c *C) {
		l.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}
