package controller

import (
	"net/http"
	"scheduler/account"
	"scheduler/router"
)

func RegisterUser(c *router.C) error {
	var u account.User
	if e := c.ParseJSON(&u); e != nil {
		return cast(e)
	}
	if e := u.Store(); e != nil {
		return e
	}
	return c.JSON(http.StatusOK, u)
}

// type userPatch struct {
// 	account.User
// 	Verify string `json:"verify"`
// }

// func (u *userPatch) Validate() error {
// 	err := u.User.Validate().(errors.FormError)
// 	if u.Verify != "" && u.Password != "" {
// 		err.Set("verify", "Need old password to update")
// 	}
// 	if err.NotOK() {
// 		return err
// 	}
// 	return nil
// }

func UpdateUser(c *router.C) error {
	// var u userPatch
	// if e := c.ParseJSON(&u); e != nil {
	// 	return cast(e)
	// }
	// return u.Patch()
	return nil

}
