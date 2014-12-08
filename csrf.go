package copter

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
	"net/http"
)

/************************************/
/************* Copter ***************/
/************************************/

func (c *Copter) EnableCSRF() {
	c.csrf = true
}

func (c *Copter) csrfHandler() http.Handler {
	if c.csrf {
		csrf := nosurf.New(c)
		csrf.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context := createContext(w, r, httprouter.Params{}, c.render)
			c.csrfHandlerFunc(context)
		}))
		return csrf
	}

	return http.Handler(c)
}

func (c *Copter) CSRFFailed(h HandlerFunc) {
	c.csrfHandlerFunc = h
}

/************************************/
/************* Context **************/
/************************************/

func (c *C) CSRFToken() string {
	return nosurf.Token(c.Request)
}
