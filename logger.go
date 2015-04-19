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
		c.Next()

		l.Printf("%s %s %v %s in %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), http.StatusText(c.Writer.Status()), time.Since(start))
	}
}
