package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type UpdateListItemService struct {
	repo domain.ListsRepository
}

func NewUpdateListItemService(repo domain.ListsRepository) *UpdateListItemService {
	return &UpdateListItemService{repo}
}

func (s *UpdateListItemService) UpdateListItem(ctx context.Context, itemID int32, listID int32, title domain.ItemTitleValueObject, description domain.ItemDescriptionValueObject, userID int32) (*domain.ListItemRecord, error) {
	foundItem, err := s.repo.FindListItem(ctx, &domain.ListItemRecord{ID: itemID, ListID: listID, UserID: userID})
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
