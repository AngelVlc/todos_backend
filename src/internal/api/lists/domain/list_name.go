package domain

import (
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type ListName string

func NewListName(name string) (ListName, error) {
	if len(name) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The list name can not be empty"}
	}

	return ListName(name), nil
}

func (l ListName) CheckIfAlreadyExists(userID int32, repo ListsRepository) error {
	existsList, err := repo.ExistsList(l, userID)
	if err != nil {
		return err
	}

	if existsList {
		return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
	}

	return nil
}
