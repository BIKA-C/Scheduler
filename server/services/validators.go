package services

import (
	"scheduler/account"
	"scheduler/router/errors"
	"scheduler/util"
)

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

func (s Account) validateAccount(a *account.Account, password account.Password) error {
	err := errors.DefaultAccError
	if a.Email == "" && a.Password == "" {
		return errors.ErrBadRequest
	}
	if a.Meta != account.EmptyMeta {
		return errors.ErrBadRequest
	}
	if a.Password != "" {
		if passwordErr := a.Password.Validate(); passwordErr != nil {
			err.Set("password", passwordErr.Error())
		}
	}
	if a.Email != "" && !util.EmailRegex.MatchString(a.Email) {
		err.Set("email", "Not a valid email")
	}
	if !s.VerifyPassword(a.UUID, password) {
		return errors.ErrUnauthorized
	}

	return nil
}
