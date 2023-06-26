package router

import (
	"log"
	"net/http"
	"os"
	"time"
)

type logger struct {
	*log.Logger
}

// NewLogger returns a new Logger instance
func Logger() HandlerFunc {
	l := &logger{log.New(os.Stdout, "[ace] ", 0)}

	return func(c *C) error {
		start := time.Now()
		e := c.Next()

		l.Printf("%s %s %v %s in %v", c.Request.Method, c.Request.URL.Path, c.writer.Status(), http.StatusText(c.writer.Status()), time.Since(start))
		return e
	}
}
