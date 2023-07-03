package controller

import (
	"scheduler/account"
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/util"
)

func (c *Controller) RegisterUser(ctx *router.C) error {
	var u account.User
	if e := ctx.ParseJSON(&u); e != nil {
		return errors.ErrInvalidJSON.Wrap(e)
	}
	if e := u.Validate(); e != nil {
		return e
	}
	if e := u.Commit(c.repo.user); e != nil {
		return e
	}
	return ctx.JSON(u)
}

func (c *Controller) UpdateUser(ctx *router.C) error {
	var u account.User
	if e := ctx.ParseJSON(&u); e != nil {
		return errors.ErrInvalidJSON.Wrap(e)
	}

	return u.Commit(c.repo.user)
}

func (c *Controller) FetchUser(ctx *router.C) error {
	u := c.repo.user.Get(util.UUID(ctx.Param("uuid")))
	if u.UUID == "" {
		return errors.BadRequest("user does not exist")
	}
	return ctx.JSON(u)
}
