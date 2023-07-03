package api

import (
	"log"
	"net/http"
	"scheduler/api/controller"
	"scheduler/router"
	"scheduler/router/errors"
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
		c.Status(http.StatusInternalServerError).JSON(errors.ErrInternalServerError)
	}
}

func (a *api) registerUserAPI() {
	user := a.router.Group("/user").Use()

	user.POST("/", a.controller.RegisterUser)                //ok
	user.PATCH("/:uuid/account", a.controller.UpdateAccount) //ok
	user.PATCH("/:uuid/", a.controller.UpdateUser)           //ok
	user.GET("/:uuid/", a.controller.FetchUser)              //ok

	viewer := a.router.Group("/user").Use(
	// todo credential check
	)

	viewer.GET("/:uuid/profile" /* todo fetchProfile */)
	viewer.PATCH("/:uuid/profile" /* todo fetchProfile */)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
