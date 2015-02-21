package ace

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

const (
	CookieSession = "cookie"
	RedisSession  = "redis"
	MongoSession  = "mongo"
)

type session struct {
	session  *sessions.Session
	isWriten bool
	isNew    bool
}

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

func (s *session) isEmpty(v interface{}) bool {
	return v == nil
}

//GetString return string value
func (s *session) GetString(key string) string {
	if s.isEmpty(s.session.Values[key]) {
		return ""
	}
	return s.session.Values[key].(string)
}

//GetInt return int value
func (s *session) GetInt(key string) int {
	if s.isEmpty(s.session.Values[key]) {
		return 0
	}
	return s.session.Values[key].(int)
}

//GetFloat64 return float64 value
func (s *session) GetFloat64(key string) float64 {
	if s.isEmpty(s.session.Values[key]) {
		return 0
	}
	return s.session.Values[key].(float64)
}

//GetBool return bool value
func (s *session) GetBool(key string) bool {
	if s.isEmpty(s.session.Values[key]) {
		return false
	}
	return s.session.Values[key].(bool)
}

//Set set session value
func (s *session) Set(key string, v interface{}) {
	s.session.Values[key] = v
	s.isWriten = true
}

func (s *session) SetInt(key string, v int) {
	s.session.Values[key] = v
	s.isWriten = true
}

func (s *session) SetFloat64(key string, v float64) {
	s.session.Values[key] = v
	s.isWriten = true
}

func (s *session) SetBool(key string, v bool) {
	s.session.Values[key] = v
	s.isWriten = true
}

//AddFlash add flash message into session
func (s *session) AddFlash(value interface{}, vars ...string) {
	s.session.AddFlash(value, vars...)
	s.isWriten = true
}

//Flashes get all flash message then remove all from session
func (s *session) Flashes(vars ...string) []interface{} {
	s.isWriten = true
	return s.session.Flashes(vars...)
}

//Delete session value
func (s *session) Delete(key string) {
	delete(s.session.Values, key)
	s.isWriten = true
}

//Clear remove all session value
func (s *session) Clear() {
	for key := range s.session.Values {
		delete(s.session.Values, key)
	}
	s.isWriten = true
}

//IsNew check session have created before
func (s *session) IsNew() bool {
	return s.isNew
}

//Session use session middleware
func (a *Ace) Session(name string, store sessions.Store, options *SessionOptions) {
	sessionStore := store

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

	a.Use(func(c *C) {
		defer context.Clear(c.Request)
		s, _ := sessionStore.Get(c.Request, name)
		if sessionOptions != nil {
			s.Options = sessionOptions
		}
		c.Session = &session{s, false, s.IsNew}
		c.Writer.Before(func(ResponseWriter) {
			if c.Session.isWriten {
				c.Session.session.Save(c.Request, c.Writer)
			}
		})
		c.Next()
	})
}
