package domain

import (
	"context"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ListNameValueObject string

func NewListNameValueObject(name string) (ListNameValueObject, error) {
	if len(name) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The list name can not be empty"}
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
