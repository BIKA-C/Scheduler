package controller

import (
	"scheduler/account"
	"scheduler/router"
	"scheduler/router/errors"
	"scheduler/util"
)

func (c *Controller) RegisterUser(ctx *router.C) error {
	var u account.User
	if err := ctx.ParseJSON(&u); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.Commit(c.repo.user); err != nil {
		return err
	}
	return ctx.JSON(u)
}

func (c *Controller) UpdateUser(ctx *router.C) error {
	var u account.User
	if err := ctx.ParseJSON(&u); err != nil {
		return errors.ErrInvalidJSON.Wrap(err)
	}
	if err := validateUserUpdate(&u); err != nil {
		return err
	}
	if err := u.Commit(c.repo.user); err != nil {
		return errors.BadRequest("User does not exist")
	}

	return ctx.JSON(u)
}

func (c *Controller) FetchUser(ctx *router.C) error {
	u := c.repo.user.Get(util.UUID(ctx.Param("uuid")))
	if u.UUID == "" {
		return errors.BadRequest("user does not exist")
	}
	return ctx.JSON(u)
}

func validateUserUpdate(u *account.User) error {
	if u.Account != account.EmptyAccount {
		return errors.Unauthorized("Can not update account")
	}
	if !u.Asset.IsEmptyAsset() {
		return errors.Unauthorized("Can not update asset")
	}
	if u.Name == "" {
		return errors.BadRequest("Name Can not be empty")
	}
	return nil
}
