package router

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"scheduler/internal/json"
	"scheduler/util"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	contentType    = "Content-Type"
	acceptLanguage = "Accept-Language"
	abortIndex     = math.MaxInt8 / 2
)

// C is context for every goroutine
type C struct {
	writerCache responseWriter
	params      httprouter.Params
	Request     *http.Request
	Writer      ResponseWriter
	index       int8
	handlers    []HandlerFunc
	data        map[string]any
}

func (a *Router) createContext(w http.ResponseWriter, r *http.Request) *C {
	c := a.pool.Get().(*C)
	c.writerCache.reset(w)
	c.Request = r
	c.index = -1
	c.data = nil

	return c
}

func (c *C) Status(status int) *C {
	c.Writer.WriteHeader(status)
	return c
}

// JSON response with application/json; charset=UTF-8 Content type
func (c *C) JSON(v any) error {
	if v == nil {
		return nil
	}

	c.Writer.Header().Set(contentType, "application/json; charset=UTF-8")

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	_, e := c.Writer.Write(buf.Bytes())
	return e
}

// String response with text/html; charset=UTF-8 Content type
func (c *C) String(format string, val ...any) error {
	c.Writer.Header().Set(contentType, "text/html; charset=UTF-8")
	if len(val) == 0 {
		c.Writer.Write(util.S2B(format))
		return nil
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	buf.WriteString(fmt.Sprintf(format, val...))

	_, e := c.Writer.Write(buf.Bytes())
	return e
}

// Write response
func (c *C) Write(v []byte) error {
	_, e := c.Writer.Write(v)
	return e
}

// Param get param from route
func (c *C) Param(name string) string {
	return c.params.ByName(name)
}

// ParseJSON decode json to any
func (c *C) ParseJSON(v any) error {
	defer c.Request.Body.Close()
	d := json.NewDecoder(c.Request.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

// Parse decode body based on application-type
func (c *C) Parse(v any) error {
	panic("not implemented")
}

// Lang get first language from HTTP Header
func (c *C) Lang() string {
	langStr := c.Request.Header.Get(acceptLanguage)
	return strings.Split(langStr, ",")[0]
}

// Redirect 302 response
func (c *C) Redirect(url string) {
	http.Redirect(c.Writer, c.Request, url, 302)
}

// Abort stop middleware
func (c *C) Abort() {
	c.index = abortIndex
}

// AbortWithStatus stop middleware and return http status code
func (c *C) AbortWithStatus(status int) {
	c.Writer.WriteHeader(status)
	c.Abort()
}

func (c *C) AbortWithError(e error) error {
	c.Abort()
	return e
}

// Next next middleware
func (c *C) Next() error {
	if c.index >= abortIndex {
		return errors.New("Aborted")
	}
	c.index++
	s := int8(len(c.handlers))
	if c.index < s {
		if err := c.handlers[c.index](c); err != nil {
			return err
		}
	}
	return nil
}

func (c *C) Method() string {
	return c.Request.Method
}

// ClientIP get ip from RemoteAddr
func (c *C) ClientIP() string {
	return c.Request.RemoteAddr
}

// Set data
func (c *C) Set(key string, v any) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data[key] = v
}

func (c *C) Delete(key string) {
	delete(c.data, key)
}

// SetAll data
func (c *C) SetAll(data map[string]any) {
	c.data = data
}

// Get data
func (c *C) Get(key string) any {
	return c.data[key]
}

// GetAllData return all data
func (c *C) GetAll() map[string]any {
	return c.data
}

func (c *C) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *C) MustQueryInt(key string, d int) int {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}

	return i
}

func (c *C) MustQueryFloat64(key string, d float64) float64 {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(c.Request.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (c *C) MustQueryString(key, d string) string {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}

	return val
}

func (c *C) MustQueryStrings(key string, d []string) []string {
	val := c.Request.URL.Query()[key]
	if len(val) == 0 {
		return d
	}

	return val
}

func (c *C) MustQueryTime(key string, layout string, d time.Time) time.Time {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, c.Request.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}

	return t
}

/////////////////////////

func (c *C) FormValue(key string) string {
	return c.Request.PostFormValue(key)
}

func (c *C) MustFormInt(key string, d int) int {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}

	return i
}

func (c *C) MustFormFloat64(key string, d float64) float64 {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(c.Request.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (c *C) MustFormString(key, d string) string {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}

	return val
}

func (c *C) MustFormStrings(key string, d []string) []string {
	if c.Request.PostForm == nil {
		c.Request.ParseForm()
	}

	val := c.Request.PostForm[key]
	if len(val) == 0 {
		return d
	}

	return val
}

func (c *C) MustFormTime(key string, layout string, d time.Time) time.Time {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, c.Request.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}

	return t
}

func (c *C) Panic(err error) {
	if err != nil {
		panic(err)
	}
}
