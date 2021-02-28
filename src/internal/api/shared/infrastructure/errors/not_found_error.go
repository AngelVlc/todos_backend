package errors

import "fmt"

// NotFoundError happens when the document does not exist in the store
type NotFoundError struct {
	Model string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v not found", e.Model)
}
