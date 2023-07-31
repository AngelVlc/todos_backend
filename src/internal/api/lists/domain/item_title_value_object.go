package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemTitleValueObject struct {
	itemTitle string
}

const itemTitleMaxLength = 50

func NewItemTitleValueObject(title string) (ItemTitleValueObject, error) {
	if len(title) == 0 {
		return ItemTitleValueObject{}, &appErrors.BadRequestError{Msg: "The item title can not be empty"}
	}

	if len(title) > itemTitleMaxLength {
		return ItemTitleValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The item title can not have more than %v characters", itemTitleMaxLength)}
	}

	return ItemTitleValueObject{itemTitle: title}, nil
}

func (v ItemTitleValueObject) String() string {
	return v.itemTitle
}

func (v ItemTitleValueObject) MarshalText() ([]byte, error) {
	return []byte(v.itemTitle), nil
}

func (v *ItemTitleValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewItemTitleValueObject(string(d))
	return err
}

func (v ItemTitleValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *ItemTitleValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewItemTitleValueObject(fmt.Sprintf("%s", sv))
		return nil

	}
	return errors.New("failed to scan ItemTitleValueObject")
}
