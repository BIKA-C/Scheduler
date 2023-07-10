package account

import (
	"scheduler/router/errors"
	"scheduler/util"
)

type Contact struct {
	Address Address `json:"address,omitempty"`
	Phone   string  `json:"phone,omitempty"`
	Email   string  `json:"email"`
}

func (c *Contact) Validate() error {
	err := c.Address.Validate().(errors.AccumulateError)

	if c.Phone != "" && !util.PhoneRegex.MatchString(c.Phone) {
		err.Set("phone", "Not a phone number")
	}

	if !util.EmailRegex.MatchString(c.Email) {
		err.Set("email", "Not a email")
	}

	return err.Build()
}
