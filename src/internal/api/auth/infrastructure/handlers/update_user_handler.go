package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	input, _ := h.RequestInput.(*infrastructure.UpdateUserInput)

	if len(input.Password) > 0 && input.Password != input.ConfirmPassword {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "Passwords don't match"}}
	}

	srv := application.NewUpdateUserService(h.UsersRepository, h.PassGen)
	user, err := srv.UpdateUser(r.Context(), userID, input.Name, input.Password, input.IsAdmin)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.UserResponse{
		ID:      user.ID,
		Name:    user.Name.String(),
		IsAdmin: user.IsAdmin,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
