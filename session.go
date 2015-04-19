package ace

import (
	"github.com/plimble/sessions"
)

//SessionOptions session options
type SessionOptions struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

//Session use session middleware
func Session(store sessions.Store, options *SessionOptions) HandlerFunc {
	var sessionOptions *sessions.Options

	if options != nil {
		sessionOptions = &sessions.Options{
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   options.MaxAge,
			Secure:   options.Secure,
			HttpOnly: options.HTTPOnly,
		}
	}

	manager := sessions.New(512, store, sessionOptions)

	return func(c *C) {
		c.sessions = manager.GetSessions(c.Request)
		defer manager.Close(c.sessions)

		c.Writer.Before(func(ResponseWriter) {
			c.sessions.Save(c.Writer)
		})

		c.Next()
	}
}
