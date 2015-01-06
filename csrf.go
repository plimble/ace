package ace

import (
	"github.com/plimble/csrf"
)

type CSRFOptions struct {
	FailedHandler HandlerFunc
}

func CSRF(options *CSRFOptions) HandlerFunc {
	if options.FailedHandler == nil {
		options.FailedHandler = defaultCSRFFailedHandler
	}

	cs := csrf.New()

	return func(c *C) {
		defer csrf.ClearContext(c.Request)

		if !cs.Check(c.Writer, c.Request) {
			options.FailedHandler(c)
			return
		}

		c.Next()
	}
}

func defaultCSRFFailedHandler(c *C) {
	c.String(500, "Invalid CSRF Token")
}
