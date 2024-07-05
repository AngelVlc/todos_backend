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

func (s *DeleteCategoryService) DeleteCategory(ctx context.Context, categoryID int32, userID int32) error {
	foundCategory, err := s.repo.FindCategory(ctx, domain.CategoryRecord{ID: categoryID, UserID: userID})
	if err != nil {
		return err
	}

	if err := s.repo.DeleteCategory(ctx, *foundCategory); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user category", InternalError: err}
	}

	return nil
}
