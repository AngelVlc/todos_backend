package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func DeleteListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewDeleteListService(h.ListsRepository)
	err := srv.DeleteList(r.Context(), listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
