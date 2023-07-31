package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UserPasswordValueObject struct {
	userPassword string
}

func NewUserPasswordValueObject(userPassword string) (UserPasswordValueObject, error) {
	if len(userPassword) == 0 {
		return UserPasswordValueObject{}, &appErrors.BadRequestError{Msg: "Password can not be empty"}
	}

	return UserPasswordValueObject{userPassword: userPassword}, nil
}

func (v UserPasswordValueObject) String() string {
	return v.userPassword
}

func (v UserPasswordValueObject) MarshalText() ([]byte, error) {
	return []byte(v.userPassword), nil
}

func (v *UserPasswordValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewUserPasswordValueObject(string(d))

	return err
}

func (v UserPasswordValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *UserPasswordValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewUserPasswordValueObject(fmt.Sprintf("%s", sv))
		return nil
	}

	return errors.New("failed to scan UserPasswordValueObject")
}
