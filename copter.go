package copter

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
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

func (c *Copter) handle(method, path string, handlers []HandlerFunc) {
	handlers = c.combineHandlers(handlers)
	c.httprouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		context := createContext(w, req, params, c.render)
		context.handlers = handlers
		context.notfoundHandlerFunc = c.notfoundHandlerFunc
		context.panicHandlerFunc = c.panicHandlerFunc
		context.Next()
	})
}

type RenderOptions render.Options

func (c *Copter) SetRenderOptions(options RenderOptions) {
	c.render = render.New(render.Options(options))
}

func (c *Copter) EnableCSRF() {
	c.csrf = true
}

func (c *Copter) CSRFFailed(h HandlerFunc) {
	c.csrfHandlerFunc = h
}

func (c *Copter) NotFound(h HandlerFunc) {
	c.notfoundHandlerFunc = h
	c.httprouter.NotFound = func(w http.ResponseWriter, r *http.Request) {
		context := createContext(w, r, httprouter.Params{}, c.render)
		h(context)
	}
}

func (c *Copter) Panic(h HandlerFunc) {
	c.panicHandlerFunc = h
	c.httprouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		context := createContext(w, r, httprouter.Params{}, c.render)
		context.Recovery = rcv
		h(context)
	}
}

func (c *Copter) GET(path string, handlers ...HandlerFunc) {
	c.handle("GET", path, handlers)
}

func (c *Copter) POST(path string, handlers ...HandlerFunc) {
	c.handle("POST", path, handlers)
}

func (c *Copter) Static(path string, root http.Dir) {
	fileServer := http.StripPrefix(path, http.FileServer(root))
	c.GET(path+"/*filepath", func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
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

func (c *Copter) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(c.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, c.handlers...)
	h = append(h, handlers...)
	return h
}
