package infrastructure

import (
	"encoding/json"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
)

type CreateUserInput struct {
	Name            domain.UserNameValueObject     `json:"name"`
	Password        domain.UserPasswordValueObject `json:"password"`
	ConfirmPassword string                         `json:"confirmPassword"`
	IsAdmin         bool                           `json:"isAdmin"`
}

func (i *CreateUserInput) UnmarshalJSON(data []byte) error {
	var realInput struct {
		Name            string `json:"name"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
		IsAdmin         bool   `json:"isAdmin"`
	}

	if err := json.Unmarshal(data, &realInput); err != nil {
		return err
	}

	nvo, err := domain.NewUserNameValueObject(realInput.Name)
	if err != nil {
		return err
	}

	pvo, err := domain.NewUserPasswordValueObject(realInput.Password)
	if err != nil {
		return err
	}

	*i = CreateUserInput{
		Name:            nvo,
		Password:        pvo,
		ConfirmPassword: realInput.ConfirmPassword,
		IsAdmin:         realInput.IsAdmin,
	}

	return nil
}
