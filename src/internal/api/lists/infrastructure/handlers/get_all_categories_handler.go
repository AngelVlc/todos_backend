package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := h.GetUserIDFromContext(r)

	srv := application.NewGetAllCategoriesService(h.CategoriesRepository)
	foundCategories, err := srv.GetAllCategories(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: foundCategories, StatusCode: http.StatusOK}
}
