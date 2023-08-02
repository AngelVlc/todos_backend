package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"

type LoginInput struct {
	UserName domain.UserNameValueObject     `json:"userName"`
	Password domain.UserPasswordValueObject `json:"password"`
}
