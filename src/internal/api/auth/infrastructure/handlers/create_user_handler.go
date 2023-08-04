package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	input, _ := h.RequestInput.(*infrastructure.CreateUserInput)

	if r.RequestURI == "/auth/createadmin" && input.Name.String() != "admin" {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "/auth/createadmin only can be used to create the admin user"}}
	}

	if input.Password.String() != input.ConfirmPassword {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "Passwords don't match"}}
	}

	srv := application.NewCreateUserService(h.UsersRepository, h.PassGen)
	newUser, err := srv.CreateUser(r.Context(), input.Name, input.Password.String(), input.IsAdmin)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: newUser, StatusCode: http.StatusCreated}
}
