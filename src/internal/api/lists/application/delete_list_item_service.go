package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type DeleteListItemService struct {
	repo domain.ListsRepository
}

func NewDeleteListItemService(repo domain.ListsRepository) *DeleteListItemService {
	return &DeleteListItemService{repo}
}

func (s *DeleteListItemService) DeleteListItem(itemID int32, listID int32, userID int32) error {
	foundList, err := s.repo.FindListItemByID(itemID, listID, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting the list item", InternalError: err}
	}

	if foundList == nil {
		return &appErrors.BadRequestError{Msg: "The list item does not exist"}
	}

	err = s.repo.DeleteListItem(itemID, listID, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the list item", InternalError: err}
	}

	return nil
}
