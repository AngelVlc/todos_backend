package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type createListItemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	createReq := createListItemRequest{}
	err := h.ParseBody(r, &createReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listTitle, err := domain.NewItemTitle(createReq.Title)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewCreateListItemService(h.ListsRepository)
	newItem, err := srv.CreateListItem(listID, listTitle, createReq.Description, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := ListItemResponse{
		ID:          newItem.ID,
		Title:       string(newItem.Title),
		Description: newItem.Description,
		ListID:      newItem.ListID,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
