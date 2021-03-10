package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type createUserRequest struct {
	UserName        *string `json:"userName"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirmPassword"`
	IsAdmin         *bool   `json:"isAdmin"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	createReq := createUserRequest{}
	err := h.ParseBody(r, &createReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	userName, err := domain.NewAuthUserName(createReq.UserName, true)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewAuthUserPassword(createReq.Password, true)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	if *createReq.Password != *createReq.ConfirmPassword {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "Passwords don't match"}}
	}

	srv := application.NewCreateUserService(h.AuthRepository, h.PassGen)
	id, err := srv.CreateUser(userName, password, createReq.IsAdmin)
	if err != nil {
		return results.ErrorResult{err}
	}

	return results.OkResult{id, http.StatusOK}
}
