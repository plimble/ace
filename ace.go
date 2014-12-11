package ace

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

type Ace struct {
	httprouter          *httprouter.Router
	handlers            []HandlerFunc
	notfoundHandlerFunc HandlerFunc
	failHandlerFunc     HandlerFunc
	pool                sync.Pool
}

type HandlerFunc func(c *C)

func New() *Ace {
	c := &Ace{}
	c.httprouter = httprouter.New()
	c.pool.New = func() interface{} {
		context := &C{}
		context.index = -1
		context.Writer = &context.writercache
		return context
	}
	return c
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

func (a *Ace) Use(middlewares ...HandlerFunc) {
	for _, handler := range middlewares {
		a.handlers = append(a.handlers, handler)
	}
}
