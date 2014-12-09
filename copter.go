package copter

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"sync"
)

type Copter struct {
	httprouter          *httprouter.Router
	render              *render.Render
	handlers            []HandlerFunc
	notfoundHandlerFunc HandlerFunc
	panicHandlerFunc    HandlerFunc
	pool                sync.Pool
}

type HandlerFunc func(c *C)
type RenderOptions render.Options

func New() *Copter {
	c := &Copter{}
	c.httprouter = httprouter.New()
	c.pool.New = func() interface{} {
		context := &C{render: render.New(render.Options{}), index: -1}
		return context
	}
	return c
}

func Default() *Copter {
	c := New()
	c.Use(Recovery())
	c.Use(Logger())
	return c
}

func (c *Copter) SetRenderOptions(options RenderOptions) {
	c.pool.New = func() interface{} {
		context := &C{render: render.New(render.Options(options)), index: -1}
		return context
	}
}

func (c *Copter) Run(addr string) {
	if err := http.ListenAndServe(addr, c); err != nil {
		panic(err)
	}
}

func (c *Copter) RunTLS(addr string, cert string, key string) {
	if err := http.ListenAndServeTLS(addr, cert, key, c); err != nil {
		panic(err)
	}
}

func (c *Copter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.httprouter.ServeHTTP(w, req)
}

func (c *Copter) Use(middlewares ...HandlerFunc) {
	for _, handler := range middlewares {
		c.handlers = append(c.handlers, handler)
	}
}
