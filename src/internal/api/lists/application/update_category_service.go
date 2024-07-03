package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UpdateCategoryService struct {
	repo domain.CategoriesRepository
}

func NewUpdateCategoryService(repo domain.CategoriesRepository) *UpdateCategoryService {
	return &UpdateCategoryService{repo}
}

func (s *UpdateCategoryService) UpdateCategory(ctx context.Context, categoryToUpdate *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	foundCategory, err := s.repo.FindCategory(ctx, domain.CategoryRecord{ID: categoryToUpdate.ID, UserID: categoryToUpdate.UserID})
	if err != nil {
		return nil, err
	}

	if foundCategory.Name != categoryToUpdate.Name.String() {
		if existsCategory, err := s.repo.ExistsCategory(ctx, domain.CategoryRecord{Name: categoryToUpdate.Name.String(), UserID: categoryToUpdate.UserID}); err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error checking if a category with the same name already exists", InternalError: err}
		} else if existsCategory {
			return nil, &appErrors.BadRequestError{Msg: "A category with the same name already exists", InternalError: nil}
		}
	}

	updatedCategory, err := s.repo.UpdateCategory(ctx, categoryToUpdate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user category", InternalError: err}
	}

	return updatedCategory, nil
}
