package account

import "scheduler/router/errors"

type Address struct {
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
