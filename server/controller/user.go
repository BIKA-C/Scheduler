package controller

import (
	"scheduler/account"
	"scheduler/repository"
	"scheduler/repository/database"
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/server/services"
)

type Controller struct {
	userServ services.User
	acctServ services.Account
}

func New(db *database.SQLite) *Controller {
	acc := repository.AccountRepo(db)
	return &Controller{
		userServ: services.NewUserServ(repository.UserRepo(&acc)),
		acctServ: services.NewAccountServ(acc),
	}
}

func (c *Controller) RegisterUser(ctx *router.C) error {
	var u account.User
	if err := ctx.ParseJSON(&u); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}
	if err := c.userServ.Register(&u); err != nil {
		return err
	}
	if user, err := c.userServ.Get(u.UUID); err != nil {
		return err
	} else {
		return ctx.JSON(user)
	}
	// negotiate contentType: html/json
	//...
}

func (c *Controller) UpdateUser(ctx *router.C) error {
	var u account.User
	if err := ctx.ParseJSON(&u); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}
	if err := c.userServ.Update(&u); err != nil {
		return err
	}
	if user, err := c.userServ.Get(u.UUID); err != nil {
		return errors.ErrInternalServerError.Wrap(err)
	} else {
		return ctx.JSON(user)
	}
}

func (c *Controller) FetchUser(ctx *router.C) error {
	u, err := c.userServ.Get(getUserID(ctx))
	if err != nil {
		return errors.BadRequest("user does not exist").Wrap(err)
	}
	return ctx.JSON(u)
}

func (c *Controller) EnrollCourse(ctx *router.C) error {
	user, inst, course := getUserID(ctx), getInstID(ctx), getCourseID(ctx)
	return nil
}

func (c *Controller) UpdateAccount(ctx *router.C) error {
	var a accountUpdate
	if err := ctx.ParseJSON(&a); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}

	a.UUID = getUserID(ctx)

	if err := c.acctServ.Update(&a.Account, a.Verify); err != nil {
		return err
	}
	if user, err := c.acctServ.Get(a.UUID); err != nil {
		return errors.ErrInternalServerError.Wrap(err)
	} else {
		return ctx.JSON(user)
	}
}
