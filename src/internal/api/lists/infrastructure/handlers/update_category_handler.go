package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	categoryID := h.ParseInt32UrlVar(r, "id")
	input, _ := h.RequestInput.(*infrastructure.CategoryInput)

	categoryEntity := input.ToCategoryEntity()
	categoryEntity.ID = categoryID

	srv := application.NewUpdateCategoryService(h.CategoriesRepository)
	updatedCategory, err := srv.UpdateCategory(r.Context(), categoryEntity)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: updatedCategory, StatusCode: http.StatusOK}
}
