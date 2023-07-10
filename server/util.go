package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scheduler/router"
	"scheduler/router/errors"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if wd != filepath.Dir(exe) {
		db = filepath.Join(filepath.Dir(exe), "/database/", "scheduler.db")
	} else {
		db = filepath.Join(filepath.Dir(wd), "/database/", "scheduler.db")
	}
}

var db string

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func forbidden(c *router.C) error {
	return errors.ErrForbidden
}

func errorHandler(c *router.C, e error) {
	switch e.(type) {
	case errors.AccumulateError:
		c.Status(http.StatusBadRequest).JSON(e)
	case errors.Error:
		err := e.(errors.Error)
		c.Status(err.Status).JSON(err)
	default:
		log.Println(e)
		c.Status(http.StatusInternalServerError).JSON(errors.ErrInternalServerError)
	}
}
