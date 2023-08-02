package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"

type UpdateUserInput struct {
	Name            domain.UserNameValueObject `json:"name"`
	Password        string                     `json:"password"`
	ConfirmPassword string                     `json:"confirmPassword"`
	IsAdmin         bool                       `json:"isAdmin"`
}
