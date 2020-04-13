package errors

import (
	"fmt"
)

// UnexpectedError is used for unexpected errors
type UnexpectedError struct {
	Msg           string
	InternalError error
}

func (e *UnexpectedError) Error() string {
	return e.Msg
}

// NotFoundError happens when the document does not exist in the store
type NotFoundError struct {
	Model string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v not found", e.Model)
}

// BadRequestError happens when an id is not valid
type BadRequestError struct {
	Msg           string
	InternalError error
}

func (e *BadRequestError) Error() string {
	return e.Msg
}

// UnauthorizedError happens when the request is unauthorized
type UnauthorizedError struct {
	Msg           string
	InternalError error
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}
