package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type ListTitle string

func NewListTitle(listTitle *string, isMandatory bool) (*ListTitle, error) {
	if isMandatory {
		if listTitle == nil {
			return nil, &appErrors.BadRequestError{Msg: "ListTitle is mandatory"}
		}

		if len(*listTitle) == 0 {
			return nil, &appErrors.BadRequestError{Msg: "ListTitle can not be empty"}
		}
	}

	return (*ListTitle)(listTitle), nil
}
