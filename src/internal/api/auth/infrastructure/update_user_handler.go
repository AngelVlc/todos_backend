package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type updateUserRequest struct {
	UserName        *string `json:"userName"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirmPassword"`
	IsAdmin         *bool   `json:"isAdmin"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	updateReq := updateUserRequest{}
	err := h.ParseBody(r, &updateReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	userName, err := domain.NewUserName(updateReq.UserName, false)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewUserPassword(updateReq.Password, false)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	if updateReq.Password != nil && *updateReq.Password != *updateReq.ConfirmPassword {
		return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: "Passwords don't match"}}
	}

	srv := application.NewUpdateUserService(h.AuthRepository, h.PassGen)
	user, err := srv.UpdateUser(&userID, userName, password, updateReq.IsAdmin)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := UserResponse{
		ID:      user.ID,
		Name:    string(user.Name),
		IsAdmin: user.IsAdmin,
	}

	return results.OkResult{Content: &res, StatusCode: http.StatusOK}
}
