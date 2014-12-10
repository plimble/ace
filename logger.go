package ace

import (
	"log"
	"net/http"
	"os"
	"time"
)

type logger struct {
	*log.Logger
}

// NewLogger returns a new Logger instance
func Logger() HandlerFunc {
	l := &logger{log.New(os.Stdout, "[ace] ", 0)}

	return func(c *C) {
		start := time.Now()
		l.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		l.Printf("Completed %v %s in %v", c.Writer.Status(), http.StatusText(c.Writer.Status()), time.Since(start))
	}
}
