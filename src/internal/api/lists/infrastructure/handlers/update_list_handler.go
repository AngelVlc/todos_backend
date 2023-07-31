package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)
	input, _ := h.RequestInput.(*infrastructure.ListInput)

	listRecord := input.ToListRecord()
	listRecord.ID = listID
	listRecord.UserID = userID
	for _, v := range listRecord.Items {
		v.ListID = listID
		v.UserID = userID
	}

	srv := application.NewUpdateListService(h.ListsRepository, h.EventBus)
	err := srv.UpdateList(r.Context(), listRecord)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: listRecord, StatusCode: http.StatusOK}
}
