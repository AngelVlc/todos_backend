package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetAllListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := h.GetUserIDFromContext(r)

	srv := application.NewGetAllListsService(h.ListsRepository)
	foundLists, err := srv.GetAllLists(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: foundLists, StatusCode: http.StatusOK}
}
