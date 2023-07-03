package router

import (
	"net/http"
	"sync"

	"scheduler/util"

	"github.com/julienschmidt/httprouter"
)

var bufPool = util.NewBufferPool()

type Router struct {
	*router
	httprouter   *httprouter.Router
	pool         sync.Pool
	panicFunc    PanicHandler
	notfoundFunc HandlerFunc
	errorHandler func(c *C, e error)
}

type PanicHandler func(c *C, rcv any) error
type HandlerFunc func(c *C) error

// New server
func New() *Router {
	a := &Router{}
	a.router = &router{
		prefix:   "/",
		handlers: nil,
		router:   a,
	}
	a.panicFunc = defaultPanic
	a.notfoundFunc = defaultNotfound
	a.errorHandler = defaultError
	a.httprouter = httprouter.New()
	a.pool.New = func() interface{} {
		c := &C{}
		c.index = -1
		c.Writer = &c.writerCache
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

// Default server white recovery and logger middleware
func Default() *Router {
	a := New()
	a.Use(Logger())
	return a
}

// // SetPoolSize of buffer
// func (a *Router) SetPoolSize(poolSize int) {
// 	bufPool = util.NewBufferPool(poolSize)
// }

func (a *Router) SetTrailingSlash(b bool) {
	a.httprouter.RedirectTrailingSlash = b
}

// Run server with specific address and port
// func (a *Router) Run(addr string) {
// 	if err := http.ListenAndServe(addr, a); err != nil {
// 		panic(err)
// 	}
// }

// // RunTLS server with specific address and port
// func (a *Router) RunTLS(addr string, cert string, key string) {
// 	if err := http.ListenAndServeTLS(addr, cert, key, a); err != nil {
// 		panic(err)
// 	}
// }

// ServeHTTP implement http.Handler
func (a *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.httprouter.ServeHTTP(w, req)
}
