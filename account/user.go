package account

import (
	e "errors"
	"scheduler/router/errors"
)

type UserAsset struct {
	Sum     int            `json:"sum"`
	Balance map[string]int `json:"balance"`
}

type User struct {
	Account `json:"account"`
	Name    string    `json:"name"`
	Asset   UserAsset `json:"asset"`
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

func (a UserAsset) IsEmptyAsset() bool {
	return a.Sum == 0 && len(a.Balance) == 0
}

func (a UserAsset) UnmarshalJSON(b []byte) error {
	return e.New("Not supported")
}
