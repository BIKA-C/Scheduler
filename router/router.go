package router

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

// var defaultPanic func(c *C, rcv any) error = nil

var defaultPanic = func(c *C, rcv interface{}) error {
	stack := debug.Stack()
	log.Printf("PANIC: %s\n%s", rcv, stack)
	c.Status(http.StatusInternalServerError).String("internal server error")
	return nil
}

var defaultNotfound = func(c *C) error {
	http.NotFoundHandler().ServeHTTP(c.Writer, c.Request)
	return nil
}

var defaultError = func(c *C, e error) {
	log.Println(e)
	if c.Writer.Written() {
		return
	}
	c.Status(http.StatusBadRequest).String(e.Error())
}

type router struct {
	handlers []HandlerFunc
	prefix   string
	router   *Router
}

// Use register middleware
func (r *router) Use(middlewares ...HandlerFunc) *router {
	for _, handler := range middlewares {
		r.handlers = append(r.handlers, handler)
	}
	return r
}

// GET handle GET method
func (r *router) GET(path string, handlers ...HandlerFunc) {
	r.Handle("GET", path, handlers)
}

// POST handle POST method
func (r *router) POST(path string, handlers ...HandlerFunc) {
	r.Handle("POST", path, handlers)
}

// PATCH handle PATCH method
func (r *router) PATCH(path string, handlers ...HandlerFunc) {
	r.Handle("PATCH", path, handlers)
}

// PUT handle PUT method
func (r *router) PUT(path string, handlers ...HandlerFunc) {
	r.Handle("PUT", path, handlers)
}

// DELETE handle DELETE method
func (r *router) DELETE(path string, handlers ...HandlerFunc) {
	r.Handle("DELETE", path, handlers)
}

// HEAD handle HEAD method
func (r *router) HEAD(path string, handlers ...HandlerFunc) {
	r.Handle("HEAD", path, handlers)
}

// OPTIONS handle OPTIONS method
func (r *router) OPTIONS(path string, handlers ...HandlerFunc) {
	r.Handle("OPTIONS", path, handlers)
}

// Group group route
func (r *router) Group(path string, handlers ...HandlerFunc) *router {
	handlers = r.combineHandlers(handlers)
	return &router{
		handlers: handlers,
		prefix:   r.path(path),
		router:   r.router,
	}
}

// NotFound call when route does not match
func (r *router) NotFound(h HandlerFunc) {
	r.router.notfoundFunc = h
}

// Error call when any handler returns an error
func (r *router) Error(h func(*C, error)) {
	r.router.errorHandler = h
}

// Panic call when panic was called
func (r *router) Panic(h PanicHandler) {
	r.router.panicFunc = h
}

//Handler convert ace.HandlerFunc to http.Handler
// func (r *Ace) HandlerFunc(h HandlerFunc) http.Handler {
// 	handlers := r.combineHandlers([]HandlerFunc{h})
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		c := r.ace.CreateContext(w, req)
// 		c.handlers = handlers
// 		c.Next()
// 		r.ace.pool.Put(c)
// 	})
// }

// Static server static file
// path is url path
// root is root directory
func (r *router) Static(path string, root http.Dir, handlers ...HandlerFunc) {
	path = r.path(path)
	fileServer := http.StripPrefix(path, http.FileServer(root))

	handlers = append(handlers, func(c *C) error {
		fileServer.ServeHTTP(c.Writer, c.Request)
		return nil
	})

	r.router.httprouter.Handle("GET", r.staticPath(path), func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		c := r.router.createContext(w, req)
		c.handlers = handlers
		if e := c.Next(); e != nil && r.router.errorHandler != nil {
			r.router.errorHandler(c, e)
		}
		r.router.pool.Put(c)
	})
}

// Handle handle with specific method
func (r *router) Handle(method, path string, handlers []HandlerFunc) {
	r.router.httprouter.Handle(method, r.path(path), func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := r.router.createContext(w, req)
		c.params = params
		c.handlers = r.combineHandlers(handlers)
		if e := c.Next(); e != nil && r.router.errorHandler != nil {
			r.router.errorHandler(c, e)
		}
		r.router.pool.Put(c)
	})
}

func (r *router) staticPath(p string) string {
	if p == "/" {
		return "/*filepath"
	}

	return concat(p, "/*filepath")
}

func (r *router) path(p string) string {
	if r.prefix == "/" {
		return p
	}

	return concat(r.prefix, p)
}

func (r *router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	aLen := len(r.handlers)
	h := make([]HandlerFunc, aLen)
	copy(h, r.handlers)
	h = append(h, handlers...)
	return h
}
