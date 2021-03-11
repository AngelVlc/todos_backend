package errors

// UnauthorizedError happens when the request is unauthorized
type UnauthorizedError struct {
	Msg           string
	InternalError error
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}
