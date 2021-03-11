package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type UserPassword string

func NewUserPassword(userPassword *string, isMandatory bool) (*UserPassword, error) {
	if isMandatory {
		if userPassword == nil {
			return nil, &appErrors.BadRequestError{Msg: "Password is mandatory"}
		}

		if len(*userPassword) == 0 {
			return nil, &appErrors.BadRequestError{Msg: "Password can not be empty"}
		}
	}

	return (*UserPassword)(userPassword), nil
}
