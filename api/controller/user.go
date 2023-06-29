package controller

import (
	"scheduler/account"
	"scheduler/router"
)

func RegisterUser(c *router.C) error {
	var u account.User
	if e := c.ParseJSON(&u); e != nil {
		return cast(e)
	}
	if e := u.Commit(); e != nil {
		return e
	}
	return c.JSON(u)
}

func UpdateUser(c *router.C) error {
	var u account.User
	if e := c.ParseJSON(&u); e != nil {
		return cast(e)
	}
	// return u.Patch()
	return nil

}

func FetchUser(c *router.C) error {
	// todo
	return nil
	// var u account.User
	// if e := c.ParseJSON(&u); e != nil {
	// 	return cast(e)
	// }
	// return u.Patch()

}
