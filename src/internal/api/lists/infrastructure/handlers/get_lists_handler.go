package handlers

import (
	"net/http"
	"strconv"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := h.GetUserIDFromContext(r)

	categoryID, _ := strconv.Atoi(r.URL.Query().Get("categoryId"))

	var categoryIDPtr *int32
	if categoryID != 0 {
		categoryIDInt32 := int32(categoryID)
		categoryIDPtr = &categoryIDInt32
	}

	srv := application.NewGetListsService(h.ListsRepository)
	foundLists, err := srv.GetLists(r.Context(), userID, categoryIDPtr)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: foundLists, StatusCode: http.StatusOK}
}
