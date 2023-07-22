package domain

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ListNameValueObject struct {
	listName string
}

const list_name_max_length = 50

func NewListNameValueObject(name string) (ListNameValueObject, error) {
	if len(name) == 0 {
		return ListNameValueObject{}, &appErrors.BadRequestError{Msg: "The list name can not be empty"}
	}

	if len(name) > list_name_max_length {
		return ListNameValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The list name can not have more than %v characters", list_name_max_length)}
	}

	return ListNameValueObject{listName: name}, nil
}

func (v ListNameValueObject) String() string {
	return v.listName
}

func (v ListNameValueObject) MarshalText() ([]byte, error) {
	return []byte(v.listName), nil
}

func (v *ListNameValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewListNameValueObject(string(d))
	return err
}

func (v ListNameValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (nvo *ListNameValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*nvo, _ = NewListNameValueObject(fmt.Sprintf("%s", sv))
		return nil

	}
	return errors.New("failed to scan ListNameValueObject")
}

func (l ListNameValueObject) CheckIfAlreadyExists(ctx context.Context, userID int32, repo ListsRepository) error {
	existsList, err := repo.ExistsList(ctx, &ListEntity{Name: l, UserID: userID})
	if err != nil {
		return err
	}

	if existsList {
		return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
	}

	return nil
}
