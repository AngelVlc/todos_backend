package errors

// BadRequestError happens when an id is not valid
type BadRequestError struct {
	Msg           string
	InternalError error
}

func (e *BadRequestError) Error() string {
	return e.Msg
}
