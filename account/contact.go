package account

import (
	"scheduler/errors"
)

type Contact struct {
	Address Address `json:"address"`
	Phone   string  `json:"phone,omitempty"`
	Email   string  `json:"email"`
}

func (c *Contact) Validate() error {
	err := c.Address.Validate().(errors.FormError)

	if c.Phone != "" && !phone.MatchString(c.Phone) {
		err.Set("phone", "Not a phone number")
	}

	if !email.MatchString(c.Email) {
		err.Set("email", "Not a email")
	}

	if !err.NotOK() {
		return err
	}

	return nil
}
