package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	itemID := helpers.ParseInt32UrlVar(r, "id")
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	updateReq, _ := h.RequestInput.(*domain.UpdateListItemInput)

	listItemTitle, err := domain.NewItemTitleValueObject(updateReq.Title)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listItemDescription, err := domain.NewItemDescriptionValueObject(updateReq.Description)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewUpdateListItemService(h.ListsRepository)
	item, err := srv.UpdateListItem(r.Context(), itemID, listID, listItemTitle, listItemDescription, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.ListItemResponse{
		ID:          item.ID,
		Title:       string(item.Title),
		Description: string(item.Description),
		ListID:      item.ListID,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
