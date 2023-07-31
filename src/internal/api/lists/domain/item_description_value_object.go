package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemDescriptionValueObject struct {
	itemDescription string
}

const itemDescriptionMaxLength = 500

func NewItemDescriptionValueObject(description string) (ItemDescriptionValueObject, error) {
	if len(description) > itemDescriptionMaxLength {
		return ItemDescriptionValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The item description can not have more than %v characters", itemDescriptionMaxLength)}
	}

	return ItemDescriptionValueObject{itemDescription: description}, nil
}

func (v ItemDescriptionValueObject) String() string {
	return v.itemDescription
}

func (v ItemDescriptionValueObject) MarshalText() ([]byte, error) {
	return []byte(v.itemDescription), nil
}

func (v *ItemDescriptionValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewItemDescriptionValueObject(string(d))
	return err
}

func (v ItemDescriptionValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *ItemDescriptionValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewItemDescriptionValueObject(fmt.Sprintf("%s", sv))
		return nil

	}
	return errors.New("failed to scan ItemDescriptionValueObject")
}
