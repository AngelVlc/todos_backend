package domain

import (
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemTitle string

func NewItemTitle(title string) (ItemTitle, error) {
	if len(title) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The item title can not be empty"}
	}

	return ItemTitle(title), nil
}
