package controller

import (
	"scheduler/account"
	"scheduler/router"
)

func (c *Controller) UpdateAccount(ctx *router.C) error {
	var a account.AccountUpdate
	if e := ctx.ParseJSON(&a); e != nil {
		return cast(e)
	}
	if e := a.Commit(c.repo.account); e != nil {
		return e
	}
	return ctx.JSON(a)
}
