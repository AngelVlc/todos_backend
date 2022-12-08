package domain

import (
	"context"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ListNameValueObject string

const list_name_max_length = 50

func NewListNameValueObject(name string) (ListNameValueObject, error) {
	if len(name) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The list name can not be empty"}
	}

	if len(name) > list_name_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The list name can not have more than %v characters", list_name_max_length)}
	}

	return ListNameValueObject(name), nil
}

func (l ListNameValueObject) CheckIfAlreadyExists(ctx context.Context, userID int32, repo ListsRepository) error {
	existsList, err := repo.ExistsList(ctx, &List{Name: l, UserID: userID})
	if err != nil {
		return err
	}

	if existsList {
		return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
	}

	return nil
}
