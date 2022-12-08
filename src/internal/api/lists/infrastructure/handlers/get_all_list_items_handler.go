package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetAllListItemsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewGetAllListItemsService(h.ListsRepository)
	foundLists, err := srv.GetAllListItems(r.Context(), listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := make([]infrastructure.ListItemResponse, len(foundLists))

	for i, v := range foundLists {
		res[i] = infrastructure.ListItemResponse{
			ID:          v.ID,
			Title:       string(v.Title),
			Description: string(v.Description),
			ListID:      v.ListID,
			Position:    v.Position,
		}
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
