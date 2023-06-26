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
		c.JSON(http.StatusBadRequest, e)
	case errors.Error:
		err := e.(errors.Error)
		c.JSON(err.Status, err)
	default:
		log.Println(e)
		c.String(http.StatusBadRequest, e.Error())
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
