package server

import (
	"fmt"
	"net/http"
	"runtime"
	"scheduler/repository/database"
	"scheduler/router"
	"scheduler/server/controller"
)

type server struct {
	api        *router.Router
	web        *router.Router
	controller *controller.Controller
}

func NewWithDB(db *database.SQLite) *server {
	r := &server{
		api:        router.New(),
		web:        router.New(),
		controller: controller.New(db),
	}
	r.setup()
	runtime.GC()
	return r
}

func New() *server {
	return NewWithDB(database.NewSQLite(db, 100, database.DefaultPragma))
}

type entry struct {
	method  string
	remark  string
	path    string
	handler router.HandlerFunc
}

func (a *server) setup() {
	a.api.NotFound(forbidden)
	a.api.Error(errorHandler)
	a.api.SetTrailingSlash(false)

	a.initUser()
}

func (a *server) initUser() {
	entries := []entry{
		{
			method:  "POST",
			path:    "/",
			handler: a.controller.RegisterUser,
		}, {
			method:  "PATCH",
			path:    fmt.Sprintf("/%s/account/", controller.UserID),
			handler: a.controller.UpdateAccount,
		}, {
			method:  "PATCH",
			path:    fmt.Sprintf("/%s/", controller.UserID),
			handler: a.controller.UpdateUser,
		}, {
			method:  "GET",
			path:    fmt.Sprintf("/%s/", controller.UserID),
			handler: a.controller.FetchUser,
		}, {
			method:  "POST",
			path:    fmt.Sprintf("/%s/", controller.UserID),
			handler: a.controller.FetchUser,
		},
	}
	user := a.api.Group("/user")
	for _, e := range entries {
		user.Handle(e.method, e.path, e.handler)
	}

}

func (a *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Host {
	case "api.localhost":
		a.api.ServeHTTP(w, r)
		break
	case "localhost":
	case "web.localhost":
		a.web.ServeHTTP(w, r)
		break
	default:
		notFound(w, r)
	}
}
