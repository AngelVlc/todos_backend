package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	createReq, _ := h.RequestInput.(*domain.CreateUserInput)

	if r.RequestURI == "/auth/createadmin" && createReq.Name != "admin" {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "/auth/createadmin only can be used to create the admin user"}}
	}

	userName, err := domain.NewUserNameValueObject(createReq.Name)
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

	srv := application.NewCreateUserService(h.UsersRepository, h.PassGen)
	newUser, err := srv.CreateUser(r.Context(), userName, password, createReq.IsAdmin)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.UserResponse{
		ID:      newUser.ID,
		Name:    string(newUser.Name),
		IsAdmin: newUser.IsAdmin,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusCreated}
}
