package domain

import (
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemTitleValueObject string

const item_title_max_length = 50

func NewItemTitleValueObject(title string) (ItemTitleValueObject, error) {
	if len(title) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The item title can not be empty"}
	}

	if len(title) > item_title_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The item title can not have more than %v characters", item_title_max_length)}
	}

	return ItemTitleValueObject(title), nil
}
