package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type AuthUserName string

func NewAuthUserName(userName *string) (*AuthUserName, error) {
	if userName == nil {
		return nil, &appErrors.BadRequestError{Msg: "UserName is mandatory"}
	}

	if len(*userName) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "UserName can not be empty"}
	}

	return (*AuthUserName)(userName), nil
}
