package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func MoveListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := h.ParseInt32UrlVar(r, "id")
	userID := h.GetUserIDFromContext(r)
	input, _ := h.RequestInput.(*infrastructure.MoveListItemInput)

	srv := application.NewMoveListItemService(h.ListsRepository, h.EventBus)
	if err := srv.MoveListItem(r.Context(), listID, input.OriginListItemID, input.DestinationListID, userID); err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: nil, StatusCode: http.StatusOK}
}
