package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResp(t *testing.T) {
	assert := assert.New(t)

	data := map[string]interface{}{
		"s": "test",
		"n": 123,
		"b": true,
	}

	a := New()
	a.GET("/", func(c *C) error {
		c.JSON(data)
		return nil
	})

	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(data)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal(buf.String(), w.Body.String())
	assert.Equal("application/json; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestStringResp(t *testing.T) {
	assert := assert.New(t)
	a := New()
	a.GET("/", func(c *C) error {
		c.String("123")
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
	assert.Equal("text/html; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestDownloadResp(t *testing.T) {
	assert := assert.New(t)
	a := New()
	a.GET("/", func(c *C) error {
		c.Write([]byte("123"))
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
}

func TestCData(t *testing.T) {
	assert := assert.New(t)
	a := New()

	a.Use(func(c *C) error {
		c.Set("test", "123")
		return c.Next()
	})

	a.GET("/", func(c *C) error {
		c.GetAll()
		c.String(c.Get("test").(string))
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
}
