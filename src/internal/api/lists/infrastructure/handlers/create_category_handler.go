package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := h.GetUserIDFromContext(r)
	input, _ := h.RequestInput.(*infrastructure.CategoryInput)

	categoryEntity := input.ToCategoryEntity()
	categoryEntity.UserID = userID

	srv := application.NewCreateCategoryService(h.CategoriesRepository)
	err := srv.CreateCategory(r.Context(), categoryEntity)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: categoryEntity, StatusCode: http.StatusCreated}
}
