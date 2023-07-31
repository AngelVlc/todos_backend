package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UserNameValueObject struct {
	userName string
}

const userNameMaxLength = 10

func NewUserNameValueObject(name string) (UserNameValueObject, error) {
	if len(name) == 0 {
		return UserNameValueObject{}, &appErrors.BadRequestError{Msg: "The user name can not be empty"}
	}

	if len(name) > userNameMaxLength {
		return UserNameValueObject{}, &appErrors.BadRequestError{Msg: fmt.Sprintf("The user name can not have more than %v characters", userNameMaxLength)}
	}

	return UserNameValueObject{userName: name}, nil
}

func (v UserNameValueObject) String() string {
	return v.userName
}

func (v UserNameValueObject) MarshalText() ([]byte, error) {
	return []byte(v.userName), nil
}

func (v *UserNameValueObject) UnmarshalText(d []byte) error {
	var err error
	*v, err = NewUserNameValueObject(string(d))

	return err
}

func (v UserNameValueObject) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *UserNameValueObject) Scan(value interface{}) error {
	if sv, err := driver.String.ConvertValue(value); err == nil {
		*v, _ = NewUserNameValueObject(fmt.Sprintf("%s", sv))
		return nil
	}

	return errors.New("failed to scan UserNameValueObject")
}
