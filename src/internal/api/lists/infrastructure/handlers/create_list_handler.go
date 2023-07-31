package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	input, _ := h.RequestInput.(*infrastructure.ListInput)

	listRecord := input.ToListRecord()
	listRecord.UserID = userID
	for _, v := range listRecord.Items {
		v.UserID = userID
	}

	srv := application.NewCreateListService(h.ListsRepository, h.EventBus)
	if err := srv.CreateList(r.Context(), listRecord); err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: listRecord, StatusCode: http.StatusCreated}
}
