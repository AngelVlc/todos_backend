package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type DeleteCategoryService struct {
	repo domain.CategoriesRepository
}

func NewDeleteCategoryService(repo domain.CategoriesRepository) *DeleteCategoryService {
	return &DeleteCategoryService{repo}
}

func (s *DeleteCategoryService) DeleteCategory(ctx context.Context, categoryID int32) error {
	foundList, err := s.repo.FindCategory(ctx, domain.CategoryEntity{ID: categoryID})
	if err != nil {
		return err
	}

	if err := s.repo.DeleteCategory(ctx, *foundList); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the category", InternalError: err}
	}

	return nil
}
