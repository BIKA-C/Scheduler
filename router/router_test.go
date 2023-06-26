package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plimble/utils/errors2"
	"github.com/stretchr/testify/assert"
)

var testHandler = func(c *C) error { return c.Next() }

func TestHTTPMethod(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	a.GET("/test", func(c *C) error {
		c.String(200, "Test")
		return nil
	})

	a.POST("/test", func(c *C) error {
		c.String(200, c.Request.FormValue("test"))
		return nil
	})

	a.PUT("/", func(c *C) error {
		c.String(200, c.Request.FormValue("test"))
		return nil
	})

	a.PATCH("/", func(c *C) error {
		c.String(200, c.Request.FormValue("test"))
		return nil
	})

	a.DELETE("/", func(c *C) error {
		c.String(200, "deleted")
		return nil
	})

	a.OPTIONS("/", func(c *C) error {
		c.String(200, "options")
		return nil
	})

	a.HEAD("/test", func(c *C) error {
		c.String(200, "head")
		return nil
	})

	r, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("Test", w.Body.String())

	r, _ = http.NewRequest("POST", "/test", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PUT", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PATCH", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("DELETE", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("deleted", w.Body.String())

	r, _ = http.NewRequest("OPTIONS", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("options", w.Body.String())

	r, _ = http.NewRequest("HEAD", "/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("head", w.Body.String())

	//tailing slash
	r, _ = http.NewRequest("GET", "/test/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(301, w.Code)

	r, _ = http.NewRequest("POST", "/test/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(307, w.Code)
}

func TestNestedGroupRoute(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	g1 := a.Group("/g1", testHandler)
	g2 := g1.Group("/g2", testHandler)
	g3 := g2.Group("/g3", testHandler)

	g3.GET("/", func(c *C) error {
		c.String(200, "g3")
		return nil
	})

	g3.GET("/test", func(c *C) error {
		c.String(200, "g3/test")
		return nil
	})

	r, _ := http.NewRequest("GET", "/g1/g2/g3/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g3", w.Body.String())

	r, _ = http.NewRequest("GET", "/g1/g2/g3/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g3/test", w.Body.String())
}

func TestGroupRoute(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	g1 := a.Group("/g1", testHandler)
	g2 := a.Group("/g2", testHandler)

	g1.GET("/", func(c *C) error {
		c.String(200, "g1")
		return nil
	})

	g1.GET("/test", func(c *C) error {
		c.String(200, "g1/test")
		return nil
	})

	g2.POST("/", func(c *C) error {
		c.String(200, "g2")
		return nil
	})

	g2.POST("/test", func(c *C) error {
		c.String(200, "g2/test")
		return nil
	})

	r, _ := http.NewRequest("GET", "/g1/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g1", w.Body.String())

	r, _ = http.NewRequest("GET", "/g1/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g1/test", w.Body.String())

	r, _ = http.NewRequest("POST", "/g2/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g2", w.Body.String())

	r, _ = http.NewRequest("POST", "/g2/test", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("g2/test", w.Body.String())
}

func TestServeStatic(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	a.Static("/assets", "./", testHandler)

	r, _ := http.NewRequest("GET", "/assets/router.go", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)

	r, _ = http.NewRequest("GET", "/assets/test.text", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(404, w.Code)
}

func TestRouteNotFound(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	a.NotFound(func(c *C) error {
		c.String(404, "test not found")
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(404, w.Code)
	assert.Equal("test not found", w.Body.String())
}

func TestPanic2(t *testing.T) {
	assert := assert.New(t)
	a := New()

	a.Panic(func(c *C, rcv interface{}) error {
		err := rcv.(errors2.Error)
		c.JSON(err.HttpStatus(), err)
		return nil
	})

	a.GET("/", func(c *C) error {
		c.Panic(errors2.NewNotFound("not found"))
		c.String(200, "123")
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(404, w.Code)
	assert.Equal("{\"message\":\"not found\"}\n", w.Body.String())
}

func TestStaticPath(t *testing.T) {
	assert := assert.New(t)

	a := New()
	path := a.staticPath("/")
	assert.Equal("/*filepath", path)

	path = a.staticPath("/public")
	assert.Equal("/public/*filepath", path)
}
