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
	return s.repo.FindCategory(ctx, domain.CategoryEntity{ID: categoryID, UserID: userID})
}
