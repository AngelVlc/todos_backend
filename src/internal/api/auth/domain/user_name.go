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

func (u UserName) CheckIfAlreadyExists(repo AuthRepository) error {
	foundUser, err := repo.FindUserByName(u)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	if foundUser != nil {
		return &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	return nil
}
