package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := h.ParseInt32UrlVar(r, "id")
	userID := h.GetUserIDFromContext(r)
	input, _ := h.RequestInput.(*infrastructure.ListInput)

	listEntity := input.ToListEntity()
	listEntity.ID = listID
	listEntity.UserID = userID
	for _, v := range listEntity.Items {
		v.ListID = listID
		v.UserID = userID
	}

	srv := application.NewUpdateListService(h.ListsRepository, h.EventBus)
	updatedList, err := srv.UpdateList(r.Context(), listEntity)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: updatedList, StatusCode: http.StatusOK}
}
