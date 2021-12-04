package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type updateListItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func UpdateListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	itemID := helpers.ParseInt32UrlVar(r, "id")
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	updateReq := updateListItemRequest{}
	err := h.ParseBody(r, &updateReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listTitle, err := domain.NewItemTitle(updateReq.Title)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewUpdateListItemService(h.ListsRepository)
	item, err := srv.UpdateListItem(r.Context(), itemID, listID, listTitle, updateReq.Description, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := ListItemResponse{
		ID:          item.ID,
		Title:       string(item.Title),
		Description: item.Description,
		ListID:      item.ListID,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
