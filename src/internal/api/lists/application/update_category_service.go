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

func (s *UpdateCategoryService) UpdateList(ctx context.Context, categoryToUpdate *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	foundCategory, err := s.repo.FindCategory(ctx, domain.CategoryEntity{ID: categoryToUpdate.ID})
	if err != nil {
		return nil, err
	}

	if foundCategory.Name != categoryToUpdate.Name {
		if existsList, err := s.repo.ExistsCategory(ctx, domain.CategoryEntity{Name: categoryToUpdate.Name}); err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error checking if a category with the same name already exists", InternalError: err}
		} else if existsList {
			return nil, &appErrors.BadRequestError{Msg: "A category with the same name already exists", InternalError: nil}
		}
	}

	updatedList, err := s.repo.UpdateCategory(ctx, categoryToUpdate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the category", InternalError: err}
	}

	return updatedList, nil
}
