package account

import (
	"scheduler/internal"
	"scheduler/router/errors"
)

type asset struct {
	Sum     uint            `json:"sum"`
	Balance map[string]uint `json:"balance"`
}

type User struct {
	Account `json:"account"`
	Name    string `json:"name"`
	Asset   asset  `json:"asset"`
}

func (u *User) Validate() error {
	var e error
	if e = u.Account.Validate(); e == nil {
		e = errors.DefaultAccError
	}
	err := e.(errors.AccumulateError)

	if u.Name == "" {
		err.Set("name", "Name empty")
	}

	return err.Build()
}

func (u *User) Commit(s internal.Saver[User]) error {
	return s.Save(u)
}