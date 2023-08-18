package infrastructure

import (
	"encoding/json"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
)

type LoginInput struct {
	UserName domain.UserNameValueObject     `json:"userName"`
	Password domain.UserPasswordValueObject `json:"password"`
}

func (i *LoginInput) UnmarshalJSON(data []byte) error {
	var realInput struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(data, &realInput); err != nil {
		return err
	}

	nvo, err := domain.NewUserNameValueObject(realInput.UserName)
	if err != nil {
		return err
	}

	pvo, err := domain.NewUserPasswordValueObject(realInput.Password)
	if err != nil {
		return err
	}

	*i = LoginInput{
		UserName: nvo,
		Password: pvo,
	}

	return nil
}
