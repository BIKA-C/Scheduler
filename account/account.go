package account

import (
	"scheduler/internal"
	"scheduler/router/errors"
	"scheduler/util"
)

// type ID string
type Password string

type Account struct {
	Password Password  `json:"password,omitempty"`
	Email    string    `json:"email"`
	UUID     util.UUID `json:"id"`
}

func (a Account) Validate() error {
	err := errors.DefaultAccError
	if len(a.Password) < 8 {
		err.Set("password", "Password too short")
	}

	if !email.MatchString(a.Email) {
		err.Set("email", "Not a valid email")
	}

	return err.Build()
}

func (a *Account) Commit(s internal.Saver[Account]) error {
	return s.Save(a)
}

type AccountUpdate struct {
	Account
	Verify string `json:"verify"`
}

func (a AccountUpdate) Validate() error {
	if a.Verify == "" {
		return errors.BadRequest("Old Password Incorrect")
	}
	// todo verify password
	return nil
}

func (p Password) MarshalJSON() ([]byte, error) {
	return nil, nil
}
