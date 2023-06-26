package account

import (
	"scheduler/errors"
	"scheduler/util"
)

type User struct {
	Account
	Name string `json:"name"`
}

func (u *User) Validate() error {
	err := errors.DefaultFormError
	if u.Name != "" {
	}

	if err.NotOK() {
		return err
	} else {
		return nil
	}
}

func (u *User) Store() error {
	u.UUID = util.RandomString(10)
	return nil
}

func (u *User) Patch() error {
	return nil
}
