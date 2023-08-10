package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type UpdateListItemsCountService struct {
	repo domain.ListsRepository
}

func NewUpdateListItemsCountService(repo domain.ListsRepository) *UpdateListItemsCountService {
	return &UpdateListItemsCountService{repo}
}

func (s *UpdateListItemsCountService) UpdateListsItemsCount(ctx context.Context, listID int32) error {
	return s.repo.UpdateListItemsCount(ctx, listID)
}
