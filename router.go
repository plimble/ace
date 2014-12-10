package ace

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *Ace) GET(path string, handlers ...HandlerFunc) {
	a.handle("GET", path, handlers)
}

func (a *Ace) POST(path string, handlers ...HandlerFunc) {
	a.handle("POST", path, handlers)
}

func (a *Ace) PATCH(path string, handlers ...HandlerFunc) {
	a.handle("PATCH", path, handlers)
}

func (a *Ace) PUT(path string, handlers ...HandlerFunc) {
	a.handle("PUT", path, handlers)
}

func (a *Ace) DELETE(path string, handlers ...HandlerFunc) {
	a.handle("DELETE", path, handlers)
}

func (a *Ace) HEAD(path string, handlers ...HandlerFunc) {
	a.handle("HEAD", path, handlers)
}

func (a *Ace) OPTIONS(path string, handlers ...HandlerFunc) {
	a.handle("OPTIONS", path, handlers)
}

func (a *Ace) NotFound(h HandlerFunc) {
	a.notfoundHandlerFunc = h
	handlers := a.combineHandlers([]HandlerFunc{h})
	a.httprouter.NotFound = func(w http.ResponseWriter, r *http.Request) {
		context := a.CreateContext(w, r)
		context.handlers = handlers
		context.Next()
	}
}

func (a *Ace) Fail(h HandlerFunc) {
	a.failHandlerFunc = h
	handlers := a.combineHandlers([]HandlerFunc{h})
	a.httprouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		context := a.CreateContext(w, r)
		context.Recovery = rcv
		context.handlers = handlers
		context.Next()
	}
}

func (a *Ace) Handler(h HandlerFunc) http.Handler {
	handlers := a.combineHandlers([]HandlerFunc{h})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := a.CreateContext(w, r)
		context.handlers = handlers
		context.Next()
	})
}

func (a *Ace) Static(path string, root http.Dir) {
	fileServer := http.StripPrefix(path, http.FileServer(root))
	a.GET(path+"/*filepath", func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (a *Ace) handle(method, path string, handlers []HandlerFunc) {
	handlers = a.combineHandlers(handlers)
	a.httprouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		context := a.CreateContext(w, req)
		context.Params = params
		context.handlers = handlers
		context.notfoundHandlerFunc = a.notfoundHandlerFunc
		context.failHandlerFunc = a.failHandlerFunc
		context.Next()
	})
}

func (a *Ace) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(a.handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, a.handlers...)
	h = append(h, handlers...)
	return h
}
