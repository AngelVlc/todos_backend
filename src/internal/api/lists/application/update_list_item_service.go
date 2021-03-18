package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UpdateListItemService struct {
	repo domain.ListsRepository
}

func NewUpdateListItemService(repo domain.ListsRepository) *UpdateListItemService {
	return &UpdateListItemService{repo}
}

func (s *UpdateListItemService) UpdateListItem(itemID int32, listID int32, title domain.ItemTitle, description string, userID int32) (*domain.ListItem, error) {
	foundItem, err := s.repo.FindListItemByID(itemID, listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting the list item", InternalError: err}
	}

	if foundItem == nil {
		return nil, &appErrors.BadRequestError{Msg: "The list item does not exist"}
	}

	foundItem.Title = title
	foundItem.Description = description

	err = s.repo.UpdateListItem(foundItem)

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the list item", InternalError: err}
	}

	return foundItem, nil
}
