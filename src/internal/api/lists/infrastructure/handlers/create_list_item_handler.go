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

func CreateListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	createReq, _ := h.RequestInput.(*domain.CreateListItemInput)

	listItemTitle, err := domain.NewItemTitleValueObject(createReq.Title)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listItemDescription, err := domain.NewItemDescriptionValueObject(createReq.Description)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewCreateListItemService(h.ListsRepository)
	newItem, err := srv.CreateListItem(r.Context(), listID, listItemTitle, listItemDescription, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	go h.EventBus.Publish("listItemCreated", listID)

	res := infrastructure.ListItemResponse{
		ID:          newItem.ID,
		Title:       string(newItem.Title),
		Description: string(newItem.Description),
		ListID:      newItem.ListID,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusCreated}
}
