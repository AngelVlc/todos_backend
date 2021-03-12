package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UserName string

func NewUserName(userName string) (UserName, error) {
	if len(userName) == 0 {
		return "", &appErrors.BadRequestError{Msg: "UserName can not be empty"}
	}

	return UserName(userName), nil
}
