package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	srv := application.NewGetUserService(h.UsersRepository)
	user, err := srv.GetUser(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.UserResponse{
		ID:      user.ID,
		Name:    string(user.Name.String()),
		IsAdmin: user.IsAdmin,
	}

	return results.OkResult{Content: &res, StatusCode: http.StatusOK}
}
