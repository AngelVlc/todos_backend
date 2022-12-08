package domain

import (
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemDescriptionValueObject string

const item_description_max_length = 500

func NewItemDescriptionValueObject(description string) (ItemDescriptionValueObject, error) {
	if len(description) > item_description_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The item description can not have more than %v characters", item_description_max_length)}
	}

	return ItemDescriptionValueObject(description), nil
}
