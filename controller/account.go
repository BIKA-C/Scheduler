package controller

import (
	"scheduler/account"
	"scheduler/router"
)

func UpdateAccount(c *router.C) error {
	var a account.AccountUpdate
	if e := c.ParseJSON(&a); e != nil {
		return cast(e)
	}
	if e := a.Patch(); e != nil {
		return e
	}
	return c.JSON(a)
}
