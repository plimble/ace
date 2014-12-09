package copter

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (c *Copter) GET(path string, handlers ...HandlerFunc) {
	c.handle("GET", path, handlers)
}

func (c *Copter) POST(path string, handlers ...HandlerFunc) {
	c.handle("POST", path, handlers)
}

func (c *Copter) PATCH(path string, handlers ...HandlerFunc) {
	c.handle("PATCH", path, handlers)
}

func (c *Copter) PUT(path string, handlers ...HandlerFunc) {
	c.handle("PUT", path, handlers)
}

func (c *Copter) DELETE(path string, handlers ...HandlerFunc) {
	c.handle("DELETE", path, handlers)
}

func (c *Copter) HEAD(path string, handlers ...HandlerFunc) {
	c.handle("HEAD", path, handlers)
}

func (c *Copter) OPTIONS(path string, handlers ...HandlerFunc) {
	c.handle("OPTIONS", path, handlers)
}

func (c *Copter) NotFound(h HandlerFunc) {
	c.notfoundHandlerFunc = h
	handlers := c.combineHandlers([]HandlerFunc{h})
	c.httprouter.NotFound = func(w http.ResponseWriter, r *http.Request) {
		context := c.CreateContext(w, r)
		context.handlers = handlers
		context.Next()
	}
}

func (c *Copter) Fail(h HandlerFunc) {
	c.failHandlerFunc = h
	handlers := c.combineHandlers([]HandlerFunc{h})
	c.httprouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		context := c.CreateContext(w, r)
		context.Recovery = rcv
		context.handlers = handlers
		context.Next()
	}
}

func (c *Copter) Handler(h HandlerFunc) http.Handler {
	handlers := c.combineHandlers([]HandlerFunc{h})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := c.CreateContext(w, r)
		context.handlers = handlers
		context.Next()
	})
}

func (c *Copter) Static(path string, root http.Dir) {
	fileServer := http.StripPrefix(path, http.FileServer(root))
	c.GET(path+"/*filepath", func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (c *Copter) handle(method, path string, handlers []HandlerFunc) {
	handlers = c.combineHandlers(handlers)
	c.httprouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		context := c.CreateContext(w, req)
		context.Params = params
		context.handlers = handlers
		context.notfoundHandlerFunc = c.notfoundHandlerFunc
		context.failHandlerFunc = c.failHandlerFunc
		context.Next()
	})
}

func (c *Copter) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(c.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, c.handlers...)
	h = append(h, handlers...)
	return h
}
