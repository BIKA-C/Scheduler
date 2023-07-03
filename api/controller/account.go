package controller

import (
	"scheduler/account"
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/util"
)

type accountUpdate struct {
	account.Account
	Verify account.Password `json:"verify,omitempty"`
}

func (c *Controller) validateAccount(a accountUpdate) error {
	if !c.repo.account.VerifyPassword(a.UUID, a.Verify) {
		return errors.BadRequest("Old Password Incorrect")
	}
	return nil
}

func (c *Controller) UpdateAccount(ctx *router.C) error {
	var a accountUpdate
	if err := ctx.ParseJSON(&a); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}

	a.UUID = util.UUID(ctx.Param("uuid"))

	if err := c.validateAccount(a); err != nil {
		return err
	}
	if err := a.Commit(c.repo.account); err != nil {
		return err
	}
	return ctx.JSON(a)
}
