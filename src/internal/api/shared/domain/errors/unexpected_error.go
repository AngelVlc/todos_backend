package errors

// UnexpectedError is used for unexpected errors
type UnexpectedError struct {
	Msg           string
	InternalError error
}

func (e *UnexpectedError) Error() string {
	return e.Msg
}
