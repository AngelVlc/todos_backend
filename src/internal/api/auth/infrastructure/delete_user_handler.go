package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	srv := application.NewDeleteUserService(h.AuthRepository)
	err := srv.DeleteUser(userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}
	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
