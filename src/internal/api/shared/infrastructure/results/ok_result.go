package results

type OkResult struct {
	Content    interface{}
	StatusCode int
}

func (r OkResult) IsError() bool {
	return false
}
