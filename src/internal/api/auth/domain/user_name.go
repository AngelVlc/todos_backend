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
	existsUser, err := repo.ExistsUser(u)
	if err != nil {
		return err
	}

	if existsUser {
		return &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	return nil
}
