package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type ListName string

func NewListName(listName *string, isMandatory bool) (*ListName, error) {
	if isMandatory {
		if listName == nil {
			return nil, &appErrors.BadRequestError{Msg: "ListName is mandatory"}
		}

		if len(*listName) == 0 {
			return nil, &appErrors.BadRequestError{Msg: "ListName can not be empty"}
		}
	}

	return (*ListName)(listName), nil
}
