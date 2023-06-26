package controller

import (
	"scheduler/errors"
)

func cast(e error) error {
	switch e.(type) {
	case errors.FormError:
	case errors.Error:
		return e
	}
	return errors.ErrInvalidJSON
}
