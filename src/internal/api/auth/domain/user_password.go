package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UserPassword string

func NewUserPassword(userPassword string) (UserPassword, error) {
	if len(userPassword) == 0 {
		return "", &appErrors.BadRequestError{Msg: "Password can not be empty"}
	}

	return UserPassword(userPassword), nil
}
