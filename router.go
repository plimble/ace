package ace

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	handlers []HandlerFunc
	prefix   string
	ace      *Ace
}

func (r *Router) Use(middlewares ...HandlerFunc) {
	for _, handler := range middlewares {
		r.handlers = append(r.handlers, handler)
	}
}

func (r *Router) GET(path string, handlers ...HandlerFunc) {
	r.Handle("GET", path, handlers)
}

func (r *Router) POST(path string, handlers ...HandlerFunc) {
	r.Handle("POST", path, handlers)
}

func (r *Router) PATCH(path string, handlers ...HandlerFunc) {
	r.Handle("PATCH", path, handlers)
}

func (r *Router) PUT(path string, handlers ...HandlerFunc) {
	r.Handle("PUT", path, handlers)
}

func (r *Router) DELETE(path string, handlers ...HandlerFunc) {
	r.Handle("DELETE", path, handlers)
}

func (r *Router) HEAD(path string, handlers ...HandlerFunc) {
	r.Handle("HEAD", path, handlers)
}

func (r *Router) OPTIONS(path string, handlers ...HandlerFunc) {
	r.Handle("OPTIONS", path, handlers)
}

func (r *Router) Group(path string, handlers ...HandlerFunc) *Router {
	handlers = r.combineHandlers(handlers)
	return &Router{
		handlers: handlers,
		prefix:   path,
		ace:      r.ace,
	}
}

//RouteNotFound call when route does not match
func (r *Router) RouteNotFound(h HandlerFunc) {
	handlers := r.combineHandlers([]HandlerFunc{h})
	r.ace.httprouter.NotFound = func(w http.ResponseWriter, req *http.Request) {
		c := r.ace.CreateContext(w, req)
		c.handlers = handlers
		c.Next()
		r.ace.pool.Put(c)
	}
}

//Panic call when panic was called
func (r *Router) Panic(h PanicHandler) {
	r.ace.httprouter.PanicHandler = func(w http.ResponseWriter, req *http.Request, rcv interface{}) {
		c := r.ace.CreateContext(w, req)
		h(c, rcv)
		r.ace.pool.Put(c)
	}
}

//Handler convert ace.HandlerFunc to http.Handler
func (r *Router) Handler(h HandlerFunc) http.Handler {
	handlers := r.combineHandlers([]HandlerFunc{h})
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := r.ace.CreateContext(w, req)
		c.handlers = handlers
		c.Next()
		r.ace.pool.Put(c)
	})
}

func (r *Router) Static(path string, root http.Dir, handlers ...HandlerFunc) {
	path = r.path(path)
	fileServer := http.StripPrefix(path, http.FileServer(root))

	handlers = append(handlers, func(c *C) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.ace.httprouter.Handle("GET", path+"/*filepath", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		c := r.ace.CreateContext(w, req)
		c.handlers = handlers
		c.Next()
		r.ace.pool.Put(c)
	})
}

func (r *Router) Handle(method, path string, handlers []HandlerFunc) {
	handlers = r.combineHandlers(handlers)
	r.ace.httprouter.Handle(method, r.path(path), func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := r.ace.CreateContext(w, req)
		c.Params = params
		c.handlers = handlers
		c.Next()
		r.ace.pool.Put(c)
	})
}

func (r *Router) path(p string) string {
	if r.prefix == "/" {
		return p
	}

	return r.prefix + p
}

func (r *Router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	aLen := len(r.handlers)
	hLen := len(handlers)
	h := make([]HandlerFunc, aLen+hLen)
	copy(h, r.handlers)
	for i := 0; i < hLen; i++ {
		h[aLen+i] = handlers[i]
	}
	return h
}
