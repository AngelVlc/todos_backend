package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	srv := application.NewGetUserService(h.AuthRepository)
	user, err := srv.GetUser(userID)
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
