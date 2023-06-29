package account

import "scheduler/router/errors"

type Address struct {
	Number   uint8  `json:"number"`
	Street   string `json:"street"`
	Unit     string `json:"unit,omitempty"`
	Province string `json:"province"`
	Country  string `json:"country"`
	PostCode string `json:"post"`
}

func (a *Address) Validate() error {
	err := errors.DefaultAccError

	// todo

	return err.Build()
}
