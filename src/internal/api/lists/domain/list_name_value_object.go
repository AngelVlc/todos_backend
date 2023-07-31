package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ListNameValueObject struct {
	listName string
}

const listNameMaxLength = 50

func NewListNameValueObject(name string) (ListNameValueObject, error) {
	if len(name) == 0 {
		return ListNameValueObject{}, &appErrors.BadRequestError{Msg: "The list name can not be empty"}
	}

	if len(name) > listNameMaxLength {
		return ListNameValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The list name can not have more than %v characters", listNameMaxLength)}
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

func (v *ListNameValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewListNameValueObject(fmt.Sprintf("%s", sv))
		return nil
	}

	return errors.New("failed to scan ListNameValueObject")
}
