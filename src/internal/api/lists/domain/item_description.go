package domain

import (
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemDescription string

const description_max_length = 500

func NewItemDescription(description string) (ItemDescription, error) {
	if len(description) > description_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The item description can not have more than %v characters", description_max_length)}
	}

	return ItemDescription(description), nil
}
