package copter

import (
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"time"
)

var notFoundPath = "/404"
var panicPath = "/500"

type Copter struct {
	httprouter          *httprouter.Router
	render              *render.Render
	handlers            []HandlerFunc
	notfoundHandlerFunc HandlerFunc
	panicHandlerFunc    HandlerFunc
	csrfHandlerFunc     HandlerFunc
	csrf                bool
}

type HandlerFunc func(c *C)
type RenderOptions render.Options

func New() *Copter {
	return &Copter{
		httprouter: httprouter.New(),
		render:     render.New(render.Options{}),
	}
}

func Default() *Copter {
	c := New()
	c.Use(Recovery())
	c.Use(Logger())
	return c
}

func (c *Copter) SetRenderOptions(options RenderOptions) {
	c.render = render.New(render.Options(options))
}

func (c *Copter) Run(addr string) {
	if err := http.ListenAndServe(addr, c.csrfHandler()); err != nil {
		panic(err)
	}
}

func (c *Copter) RunAndGracefulShutdown(addr string, timeout time.Duration) {
	graceful.Run(addr, timeout, c.csrfHandler())
}

func (c *Copter) RunTLS(addr string, cert string, key string) {
	if err := http.ListenAndServeTLS(addr, cert, key, c.csrfHandler()); err != nil {
		panic(err)
	}
}

func (c *Copter) RunTLSAndGracefulShutdown(addr string, cert string, key string, timeout time.Duration) {
	srv := &http.Server{
		Addr:    addr,
		Handler: c.csrfHandler(),
	}

	if err := graceful.ListenAndServeTLS(srv, cert, key, timeout); err != nil {
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
