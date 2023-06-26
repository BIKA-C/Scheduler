package errors

import (
	"net/http"
)

var ErrInvalidJSON = Error{
	Status: http.StatusBadRequest,
	Title:  "Invalid JSON",
}
