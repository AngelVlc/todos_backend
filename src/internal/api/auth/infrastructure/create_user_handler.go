package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type createUserRequest struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	IsAdmin         bool   `json:"isAdmin"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	createReq := createUserRequest{}
	err := h.ParseBody(r, &createReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	if r.RequestURI == "/auth/createadmin" && createReq.Name != "admin" {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "/auth/createadmin only can be used to create the admin user"}}
	}

	userName, err := domain.NewUserName(createReq.Name)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewUserPassword(createReq.Password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	if createReq.Password != createReq.ConfirmPassword {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "Passwords don't match"}}
	}

	srv := application.NewCreateUserService(h.AuthRepository, h.PassGen)
	newUser, err := srv.CreateUser(r.Context(), userName, password, createReq.IsAdmin)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := UserResponse{
		ID:      newUser.ID,
		Name:    string(newUser.Name),
		IsAdmin: newUser.IsAdmin,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
