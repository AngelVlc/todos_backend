package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type ListName string

func NewListName(name string) (ListName, error) {
	if len(name) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The list name can not be empty"}
	}

	return ListName(name), nil
}
