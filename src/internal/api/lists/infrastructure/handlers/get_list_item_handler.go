package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
)

func GetListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "listId")
	itemID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewGetListItemService(h.ListsRepository)
	foundList, err := srv.GetListItem(r.Context(), itemID, listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.ListItemResponse{
		ID:          foundList.ID,
		Title:       string(foundList.Title),
		Description: foundList.Description,
		ListID:      foundList.ListID,
	}

	return results.OkResult{Content: &res, StatusCode: http.StatusOK}
}
