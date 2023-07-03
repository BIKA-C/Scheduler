package controller

import (
	"scheduler/account"
	"scheduler/router"
	"scheduler/router/errors"
)

func (c *Controller) UpdateAccount(ctx *router.C) error {
	var a account.AccountUpdate
	if err := ctx.ParseJSON(&a); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}
	if err := a.Validate(); err != nil {
		return err
	}
	if err := a.Commit(c.repo.account); err != nil {
		return err
	}
	return ctx.JSON(a)
}
