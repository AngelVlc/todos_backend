package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	categoryID := h.ParseInt32UrlVar(r, "id")
	userID := h.GetUserIDFromContext(r)

	srv := application.NewDeleteCategoryService(h.CategoriesRepository)
	err := srv.DeleteCategory(r.Context(), categoryID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
