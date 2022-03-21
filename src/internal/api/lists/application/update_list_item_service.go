package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
)

type UpdateListItemService struct {
	repo domain.ListsRepository
}

func NewUpdateListItemService(repo domain.ListsRepository) *UpdateListItemService {
	return &UpdateListItemService{repo}
}

func (s *UpdateListItemService) UpdateListItem(ctx context.Context, itemID int32, listID int32, title domain.ItemTitle, description string, userID int32) (*domain.ListItem, error) {
	foundItem, err := s.repo.FindListItemByID(ctx, itemID, listID, userID)
	if err != nil {
		return nil, err
	}

	foundItem.Title = title
	foundItem.Description = description

	err = s.repo.UpdateListItem(ctx, foundItem)

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the list item", InternalError: err}
	}

	return foundItem, nil
}
