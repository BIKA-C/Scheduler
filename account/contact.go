package account

import "scheduler/router/errors"

type Contact struct {
	Address Address `json:"address,omitempty"`
	Phone   string  `json:"phone,omitempty"`
	Email   string  `json:"email"`
}

func (c *Contact) Validate() error {
	err := c.Address.Validate().(errors.AccumulateError)

	if c.Phone != "" && !phone.MatchString(c.Phone) {
		err.Set("phone", "Not a phone number")
	}

	if !email.MatchString(c.Email) {
		err.Set("email", "Not a email")
	}

	return err.Build()
}
