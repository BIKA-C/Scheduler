package account

import (
	"scheduler/router/errors"
	"scheduler/util"
)

type asset struct {
	Sum     uint            `json:"sum"`
	Balance map[string]uint `json:"balance"`
}

type User struct {
	Account
	Name  string `json:"name"`
	Asset asset  `json:"asset"`
}

func (u *User) Validate() error {
	err := u.Account.Validate().(errors.AccumulateError)
	if u.Name != "" {
	}

	return err.Build()
}

func (u *User) Commit() error {
	u.UUID = util.RandomString(10)
	return nil
}

// func (u *User) Patch() error {
// 	return nil
// }
