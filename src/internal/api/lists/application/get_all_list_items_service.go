package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetAllListItemsService struct {
	repo domain.ListsRepository
}

func NewGetAllListItemsService(repo domain.ListsRepository) *GetAllListItemsService {
	return &GetAllListItemsService{repo}
}

func (s *GetAllListItemsService) GetAllListItems(ctx context.Context, listID int32, userID int32) ([]domain.ListItemRecord, error) {
	foundItems, err := s.repo.GetAllListItems(ctx, listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all list items", InternalError: err}
	}

	return foundItems, nil
}
