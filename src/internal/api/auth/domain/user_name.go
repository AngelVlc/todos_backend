package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type UserName string

func NewUserName(userName *string, isMandatory bool) (*UserName, error) {
	if isMandatory {
		if userName == nil {
			return nil, &appErrors.BadRequestError{Msg: "UserName is mandatory"}
		}

		if len(*userName) == 0 {
			return nil, &appErrors.BadRequestError{Msg: "UserName can not be empty"}
		}
	}

	return (*UserName)(userName), nil
}
