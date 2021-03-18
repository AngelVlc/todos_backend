package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type GetListItemService struct {
	repo domain.ListsRepository
}

func NewGetListItemService(repo domain.ListsRepository) *GetListItemService {
	return &GetListItemService{repo}
}

func (s *GetListItemService) GetListItem(itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	foundList, err := s.repo.FindListItemByID(itemID, listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting the list item", InternalError: err}
	}

	if foundList == nil {
		return nil, &appErrors.BadRequestError{Msg: "The list item does not exist"}
	}

	return foundList, nil
}
