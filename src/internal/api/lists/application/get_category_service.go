package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type GetCategoryService struct {
	repo domain.CategoriesRepository
}

func NewGetCategoryService(repo domain.CategoriesRepository) *GetCategoryService {
	return &GetCategoryService{repo}
}

func (s *GetCategoryService) GetCategory(ctx context.Context, categoryID int32, userID int32) (*domain.CategoryEntity, error) {
	foundCategory, err := s.repo.FindCategory(ctx, domain.CategoryRecord{ID: categoryID, UserID: userID})
	if err != nil {
		return nil, err
	}

	return foundCategory.ToCategoryEntity(), nil
}
