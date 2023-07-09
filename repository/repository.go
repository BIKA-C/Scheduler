package repository

import (
	"scheduler/router/errors"
	"scheduler/util"
)

type ID interface {
	util.ID | util.UUID | ~int | string
}

type Repository[T any, U ID] interface {
	Save(*T) error
	Get(U) (T, error)
	Delete(U) error
}

var (
	ErrItemDoesNotExist    = errors.NotFound("Item does not exist")
	ErrItemAlreadyExist    = errors.Conflict("Item already exist")
	ErrItemCanNotBeCreated = errors.InternalServerError("Item can not be created")
	ErrItemCanNotBeDeleted = errors.InternalServerError("Item can not be deleted")
	ErrItemCanNotBeUpdated = errors.InternalServerError("Item can not be updated")

	ErrAccountDoesNotExist    = errors.NotFound("Account does not exist")
	ErrAccountAlreadyExist    = errors.Conflict("Account already exist")
	ErrAccountCanNotBeCreated = errors.InternalServerError("Account can not be created")
	ErrAccountCanNotBeDeleted = errors.InternalServerError("Account can not be deleted")
	ErrAccountCanNotBeUpdated = errors.InternalServerError("Account can not be updated")

	ErrUserDoesNotExist    = errors.NotFound("User does not exist")
	ErrUserAlreadyExist    = errors.Conflict("User already exist")
	ErrUserCanNotBeCreated = errors.InternalServerError("User can not be created")
	ErrUserCanNotBeDeleted = errors.InternalServerError("User can not be deleted")
	ErrUserCanNotBeUpdated = errors.InternalServerError("User can not be updated")
)
