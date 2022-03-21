package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	srv := application.NewDeleteUserService(h.AuthRepository)
	err := srv.DeleteUser(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
