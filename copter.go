package copter

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type Copter struct {
	httprouter *httprouter.Router
	render     *render.Render
	handlers   []HandlerFunc
}

type HandlerFunc func(c *C)

func New() *Copter {
	return &Copter{
		httprouter: httprouter.New(),
		render:     render.New(render.Options{}),
	}
}

func (c *Copter) handle(method, path string, handlers []HandlerFunc) {
	handlers = c.combineHandlers(handlers)
	c.httprouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		context := &C{
			Params:   params,
			Req:      req,
			Res:      w,
			render:   c.render,
			index:    -1,
			handlers: handlers,
		}

		context.Next()
	})
}

type RenderOptions render.Options

func (c *Copter) SetRenderOptions(options RenderOptions) {
	c.render = render.New(render.Options(options))
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
		fileServer.ServeHTTP(c.Res, c.Req)
	})
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

func (c *Copter) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(c.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, c.handlers...)
	h = append(h, handlers...)
	return h
}
