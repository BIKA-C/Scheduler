package api

import (
	"log"
	"net/http"
	"scheduler/controller"
	"scheduler/errors"
	"scheduler/router"
)

type api router.Router

func New() *router.Router {
	r := (*api)(router.New())
	r.registerAccountAPI()
	r.Error(errorHandler)
	(*router.Router)(r).SetTrailingSlash(false)
	return (*router.Router)(r)
}

func errorHandler(c *router.C, e error) {
	switch e.(type) {
	case errors.FormError:
		c.Status(http.StatusBadRequest).JSON(e)
	case errors.Error:
		err := e.(errors.Error)
		c.Status(err.Status).JSON(err)
	default:
		log.Println(e)
		c.Status(http.StatusBadRequest).String(e.Error())
	}
}

func (a *api) registerAccountAPI() {
	acc := a.Group("/account")
	acc.Use(
	// todo
	)

	user := acc.Group("/user")
	user.POST("/", controller.RegisterUser)
	user.PATCH("/:uuid", controller.UpdateAccount)
}
