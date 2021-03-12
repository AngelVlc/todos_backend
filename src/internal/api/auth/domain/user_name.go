package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UserName string

func NewUserName(userName *string) (*UserName, error) {
	if userName == nil {
		return nil, &appErrors.BadRequestError{Msg: "UserName is mandatory"}
	}

	if len(*userName) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "UserName can not be empty"}
	}

	return (*UserName)(userName), nil
}
