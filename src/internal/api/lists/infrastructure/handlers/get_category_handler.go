package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetCategoryHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	categoryID := h.ParseInt32UrlVar(r, "id")
	userID := h.GetUserIDFromContext(r)

	srv := application.NewGetCategoryService(h.CategoriesRepository)
	foundCategory, err := srv.GetCategory(r.Context(), categoryID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: foundCategory, StatusCode: http.StatusOK}
}
