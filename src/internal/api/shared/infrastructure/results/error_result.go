package results

type ErrorResult struct {
	Err error
}

func (e ErrorResult) IsError() bool {
	return true
}
