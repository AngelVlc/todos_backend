package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type DeleteListItemService struct {
	repo domain.ListsRepository
}

func NewDeleteListItemService(repo domain.ListsRepository) *DeleteListItemService {
	return &DeleteListItemService{repo}
}

func (s *DeleteListItemService) DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error {
	_, err := s.repo.FindListItem(ctx, &domain.ListItemEntity{ID: itemID, ListID: listID, UserID: userID})
	if err != nil {
		return err
	}

	err = s.repo.DeleteListItem(ctx, itemID, listID, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the list item", InternalError: err}
	}

	return nil
}
