package application

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
)

type GetListItemService struct {
	repo domain.ListsRepository
}

func NewGetListItemService(repo domain.ListsRepository) *GetListItemService {
	return &GetListItemService{repo}
}

func (s *GetListItemService) GetListItem(ctx context.Context, itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	return s.repo.FindListItemByID(ctx, itemID, listID, userID)
}
