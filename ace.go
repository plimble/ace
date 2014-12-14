package ace

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

type Ace struct {
	*Router
	httprouter *httprouter.Router
	pool       sync.Pool
}

type HandlerFunc func(c *C)
type ErrorHandlerFunc func(c *C, err error)

func New() *Ace {
	a := &Ace{}
	a.Router = &Router{
		handlers: nil,
		prefix:   "/",
		ace:      a,
	}
	a.httprouter = httprouter.New()
	a.pool.New = func() interface{} {
		c := &C{}
		c.index = -1
		c.Writer = &c.writercache
		return c
	}
	return a
}

func Default() *Ace {
	a := New()
	a.Use(Recovery())
	a.Use(Logger())
	return a
}

func (a *Ace) Run(addr string) {
	if err := http.ListenAndServe(addr, a); err != nil {
		panic(err)
	}
}

func (a *Ace) RunTLS(addr string, cert string, key string) {
	if err := http.ListenAndServeTLS(addr, cert, key, a); err != nil {
		panic(err)
	}
}

func (a *Ace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.httprouter.ServeHTTP(w, req)
}
