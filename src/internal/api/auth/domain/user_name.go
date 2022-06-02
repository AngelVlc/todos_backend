package domain

import (
	"context"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UserName string

func NewUserName(userName string) (UserName, error) {
	if len(userName) == 0 {
		return "", &appErrors.BadRequestError{Msg: "UserName can not be empty"}
	}

	return UserName(userName), nil
}

func (u UserName) CheckIfAlreadyExists(ctx context.Context, repo UsersRepository) error {
	foundUser, err := repo.FindUser(ctx, &User{Name: u})
	if err != nil {
		return err
	}

	if foundUser != nil {
		return &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	return nil
}
