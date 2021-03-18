package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func DeleteListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	itemID := helpers.ParseInt32UrlVar(r, "id")
	listID := helpers.ParseInt32UrlVar(r, "listId")
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewDeleteListItemService(h.ListsRepository)
	err := srv.DeleteListItem(itemID, listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}
	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
