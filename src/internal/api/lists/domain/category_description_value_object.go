package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CategoryDescriptionValueObject struct {
	categoryDescription string
}

const categoryDescriptionMaxLength = 500

func NewCategoryDescriptionValueObject(description string) (CategoryDescriptionValueObject, error) {
	if len(description) > categoryDescriptionMaxLength {
		return CategoryDescriptionValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The category description can not have more than %v characters", categoryDescriptionMaxLength)}
	}

	return CategoryDescriptionValueObject{categoryDescription: description}, nil
}

func (v CategoryDescriptionValueObject) String() string {
	return v.categoryDescription
}

func (v CategoryDescriptionValueObject) MarshalText() ([]byte, error) {
	return []byte(v.categoryDescription), nil
}

func (v *CategoryDescriptionValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewCategoryDescriptionValueObject(string(d))
	return err
}

func (v CategoryDescriptionValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *CategoryDescriptionValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewCategoryDescriptionValueObject(fmt.Sprintf("%s", sv))
		return nil

	}
	return errors.New("failed to scan CategoryDescriptionValueObject")
}
