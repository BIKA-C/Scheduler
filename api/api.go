package api

import (
	"log"
	"net/http"
	"scheduler/api/controller"
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/router/middleware"
)

type api struct {
	router     *router.Router
	controller *controller.Controller
}

func New() *api {
	r := &api{
		router:     router.New(),
		controller: controller.New(),
	}
	r.registerUserAPI()
	r.router.Error(errorHandler)
	r.router.SetTrailingSlash(false)
	return r
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
		c.Status(http.StatusBadRequest).String(e.Error())
	}
}

func (a *api) registerUserAPI() {
	user := a.router.Group("/user").Use(
		middleware.Validate(func(c *router.C) bool {
			return c.Method() != http.MethodPost
		}),
	)

	user.POST("/", controller.RegisterUser)
	user.PATCH("/:uuid/account", a.controller.UpdateAccount)
	user.PATCH("/:uuid/", controller.UpdateUser)
	user.GET("/:uuid/", controller.FetchUser)

	viewer := a.router.Group("/user").Use(
	// todo credential check
	)

	viewer.GET("/profile" /* todo fetchProfile */)
	viewer.PATCH("/profile" /* todo fetchProfile */)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
