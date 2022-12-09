package domain

import (
	"context"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UserNameValueObject string

const user_name_max_length = 10

func NewUserNameValueObject(name string) (UserNameValueObject, error) {
	if len(name) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The user name can not be empty"}
	}

	if len(name) > user_name_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The user name can not have more than %v characters", user_name_max_length)}
	}

	return UserNameValueObject(name), nil
}

func (u UserNameValueObject) CheckIfAlreadyExists(ctx context.Context, repo UsersRepository) error {
	existsUser, err := repo.ExistsUser(ctx, &UserEntity{Name: u})
	if err != nil {
		return err
	}

	if existsUser {
		return &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	return nil
}
