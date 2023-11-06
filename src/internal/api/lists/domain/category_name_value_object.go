package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CategoryNameValueObject struct {
	categoryName string
}

const categoryNameMaxLength = 12

func NewCategoryNameValueObject(name string) (CategoryNameValueObject, error) {
	if len(name) == 0 {
		return CategoryNameValueObject{}, &appErrors.BadRequestError{Msg: "The category name can not be empty"}
	}

	if len(name) > categoryNameMaxLength {
		return CategoryNameValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The category name can not have more than %v characters", categoryNameMaxLength)}
	}

	return CategoryNameValueObject{categoryName: name}, nil
}

func (v CategoryNameValueObject) String() string {
	return v.categoryName
}

func (v CategoryNameValueObject) MarshalText() ([]byte, error) {
	return []byte(v.categoryName), nil
}

func (v *CategoryNameValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewCategoryNameValueObject(string(d))
	return err
}

func (v CategoryNameValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *CategoryNameValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewCategoryNameValueObject(fmt.Sprintf("%s", sv))
		return nil

	}
	return errors.New("failed to scan CategoryNameValueObject")
}
