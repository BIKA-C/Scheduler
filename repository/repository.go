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

	ErrInstructorDoesNotExist    = errors.NotFound("Instructor does not exist")
	ErrInstructorAlreadyExist    = errors.Conflict("Instructor already exist")
	ErrInstructorCanNotBeCreated = errors.InternalServerError("Instructor can not be created")
	ErrInstructorCanNotBeDeleted = errors.InternalServerError("Instructor can not be deleted")
	ErrInstructorCanNotBeUpdated = errors.InternalServerError("Instructor can not be updated")

	ErrCourseDoesNotExist    = errors.NotFound("Course does not exist")
	ErrCourseAlreadyExist    = errors.Conflict("Course already exist")
	ErrCourseCanNotBeCreated = errors.InternalServerError("Course can not be created")
	ErrCourseCanNotBeDeleted = errors.InternalServerError("Course can not be deleted")
	ErrCourseCanNotBeUpdated = errors.InternalServerError("Course can not be updated")

	ErrSectionDoesNotExist    = errors.NotFound("Section does not exist")
	ErrSectionAlreadyExist    = errors.Conflict("Section already exist")
	ErrSectionCanNotBeCreated = errors.InternalServerError("Section can not be created")
	ErrSectionCanNotBeDeleted = errors.InternalServerError("Section can not be deleted")
	ErrSectionCanNotBeUpdated = errors.InternalServerError("Section can not be updated")

	ErrClassDoesNotExist     = errors.NotFound("Class does not exist")
	ErrClassAlreadyExist     = errors.Conflict("Class already exist")
	ErrClassCanNotBeCreated  = errors.InternalServerError("Class can not be created")
	ErrClassCanNotBeDeleted  = errors.InternalServerError("Class can not be deleted")
	ErrClassCanNotBeUpdated  = errors.InternalServerError("Class can not be updated")
	ErrClassUpdatedForbidden = errors.Forbidden("Class update not allowed")
)
