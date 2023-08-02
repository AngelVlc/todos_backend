package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"

type CreateUserInput struct {
	Name            domain.UserNameValueObject     `json:"name"`
	Password        domain.UserPasswordValueObject `json:"password"`
	ConfirmPassword string                         `json:"confirmPassword"`
	IsAdmin         bool                           `json:"isAdmin"`
}
