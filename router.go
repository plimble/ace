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

func (c *Copter) Static(path string, root http.Dir) {
	fileServer := http.StripPrefix(path, http.FileServer(root))
	c.GET(path+"/*filepath", func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (c *Copter) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(c.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, c.handlers...)
	h = append(h, handlers...)
	return h
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
