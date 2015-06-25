package ace

import (
	"github.com/julienschmidt/httprouter"
	"github.com/plimble/utils/pool"
	"net/http"
	"sync"
)

var bufPool = pool.NewBufferPool(100)

type Ace struct {
	*Router
	httprouter   *httprouter.Router
	pool         sync.Pool
	render       Renderer
	panicFunc    PanicHandler
	notfoundFunc HandlerFunc
}

type PanicHandler func(c *C, rcv interface{})
type HandlerFunc func(c *C)

func GetPool() *pool.BufferPool {
	return bufPool
}

//New server
func New() *Ace {
	a := &Ace{}
	a.Router = &Router{
		handlers: nil,
		prefix:   "/",
		ace:      a,
	}
	a.panicFunc = defaultPanic
	a.notfoundFunc = defaultNotfound
	a.httprouter = httprouter.New()
	a.pool.New = func() interface{} {
		c := &C{}
		c.index = -1
		c.Writer = &c.writercache
		return c
	}

	a.httprouter.PanicHandler = func(w http.ResponseWriter, req *http.Request, rcv interface{}) {
		c := a.createContext(w, req)
		a.panicFunc(c, rcv)
		a.pool.Put(c)
	}

	a.httprouter.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c := a.createContext(w, req)
		a.notfoundFunc(c)
		a.pool.Put(c)
	})

	return a
}

//Default server white recovery and logger middleware
func Default() *Ace {
	a := New()
	a.Use(Logger())
	return a
}

//SetPoolSize of buffer
func (a *Ace) SetPoolSize(poolSize int) {
	bufPool = pool.NewBufferPool(poolSize)
}

//Run server with specific address and port
func (a *Ace) Run(addr string) {
	if err := http.ListenAndServe(addr, a); err != nil {
		panic(err)
	}
}

//RunTLS server with specific address and port
func (a *Ace) RunTLS(addr string, cert string, key string) {
	if err := http.ListenAndServeTLS(addr, cert, key, a); err != nil {
		panic(err)
	}
}

//ServeHTTP implement http.Handler
func (a *Ace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.httprouter.ServeHTTP(w, req)
}
