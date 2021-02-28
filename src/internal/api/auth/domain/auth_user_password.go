package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type AuthUserPassword string

func NewAuthUserPassword(userPassword *string) (*AuthUserPassword, error) {
	if userPassword == nil {
		return nil, &appErrors.BadRequestError{Msg: "Password is mandatory"}
	}

	if len(*userPassword) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "Password can not be empty"}
	}

	return (*AuthUserPassword)(userPassword), nil
}
