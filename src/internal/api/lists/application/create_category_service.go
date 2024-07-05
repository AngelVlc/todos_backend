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

func (s *CreateCategoryService) CreateCategory(ctx context.Context, categoryToCreate *domain.CategoryEntity) error {
	if existsCategory, err := s.repo.ExistsCategory(ctx, domain.CategoryRecord{Name: categoryToCreate.Name.String(), UserID: categoryToCreate.UserID}); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error checking if a category with the same name already exists", InternalError: err}
	} else if existsCategory {
		return &appErrors.BadRequestError{Msg: "A category with the same name already exists", InternalError: nil}
	}

	record := categoryToCreate.ToCategoryRecord()

	err := s.repo.CreateCategory(ctx, record)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error creating the user category", InternalError: err}
	}

	categoryToCreate.ID = record.ID

	return nil
}
