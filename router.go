package ace

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *Ace) GET(path string, handlers ...HandlerFunc) {
	a.Handle("GET", path, handlers)
}

func (a *Ace) POST(path string, handlers ...HandlerFunc) {
	a.Handle("POST", path, handlers)
}

func (a *Ace) PATCH(path string, handlers ...HandlerFunc) {
	a.Handle("PATCH", path, handlers)
}

func (a *Ace) PUT(path string, handlers ...HandlerFunc) {
	a.Handle("PUT", path, handlers)
}

func (a *Ace) DELETE(path string, handlers ...HandlerFunc) {
	a.Handle("DELETE", path, handlers)
}

func (a *Ace) HEAD(path string, handlers ...HandlerFunc) {
	a.Handle("HEAD", path, handlers)
}

func (a *Ace) OPTIONS(path string, handlers ...HandlerFunc) {
	a.Handle("OPTIONS", path, handlers)
}

func (a *Ace) RouteNotFound(h HandlerFunc) {
	handlers := a.combineHandlers([]HandlerFunc{h})
	a.httprouter.NotFound = func(w http.ResponseWriter, r *http.Request) {
		c := a.CreateContext(w, r)
		c.handlers = handlers
		c.Next()
		a.pool.Put(c)
	}
}

func (a *Ace) Error(h ErrorHandlerFunc) {
	a.errorHandlerFunc = h
}

func (a *Ace) Panic(h HandlerFunc) {
	handlers := a.combineHandlers([]HandlerFunc{h})
	a.httprouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		c := a.CreateContext(w, r)
		c.handlers = handlers
		c.Next()
		a.pool.Put(c)
	}
}

func (a *Ace) Handler(h HandlerFunc) http.Handler {
	handlers := a.combineHandlers([]HandlerFunc{h})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := a.CreateContext(w, r)
		c.handlers = handlers
		c.Next()
		a.pool.Put(c)
	})
}

func (a *Ace) Static(path string, root http.Dir) {
	fileServer := http.StripPrefix(path, http.FileServer(root))
	a.GET(path+"/*filepath", func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (a *Ace) Handle(method, path string, handlers []HandlerFunc) {
	handlers = a.combineHandlers(handlers)
	a.httprouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := a.CreateContext(w, req)
		c.Params = params
		c.handlers = handlers
		c.errorHandlerFunc = a.errorHandlerFunc
		c.Next()
		a.pool.Put(c)
	})
}

func (a *Ace) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(a.handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, a.handlers...)
	return append(mergedHandlers, handlers...)

	// aLen := len(a.handlers)
	// hLen := len(handlers)
	// h := make([]HandlerFunc, aLen+hLen)
	// copy(h, a.handlers)
	// for i := 0; i < hLen; i++ {
	// 	h[aLen+i] = handlers[i]
	// }
	// return h
}
