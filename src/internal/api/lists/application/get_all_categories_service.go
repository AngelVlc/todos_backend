package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetAllCategoriesService struct {
	repo domain.CategoriesRepository
}

func NewGetAllCategoriesService(repo domain.CategoriesRepository) *GetAllCategoriesService {
	return &GetAllCategoriesService{repo}
}

func (s *GetAllCategoriesService) GetAllCategories(ctx context.Context, userID int32) ([]*domain.CategoryEntity, error) {
	foundCategories, err := s.repo.GetAllCategoriesForUser(ctx, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all user categories", InternalError: err}
	}

	return foundCategories, nil
}
