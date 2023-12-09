package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CreateCategoryService struct {
	repo domain.CategoriesRepository
}

func NewCreateCategoryService(repo domain.CategoriesRepository) *CreateCategoryService {
	return &CreateCategoryService{repo}
}

func (s *CreateCategoryService) CreateCategory(ctx context.Context, categoryToCreate *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	if existsCategory, err := s.repo.ExistsCategory(ctx, domain.CategoryEntity{Name: categoryToCreate.Name, UserID: categoryToCreate.UserID}); err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error checking if a category with the same name already exists", InternalError: err}
	} else if existsCategory {
		return nil, &appErrors.BadRequestError{Msg: "A category with the same name already exists", InternalError: nil}
	}

	createdList, err := s.repo.CreateCategory(ctx, categoryToCreate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user category", InternalError: err}
	}

	return createdList, nil
}
