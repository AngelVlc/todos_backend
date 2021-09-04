package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type CreateListItemService struct {
	repo domain.ListsRepository
}

func NewCreateListItemService(repo domain.ListsRepository) *CreateListItemService {
	return &CreateListItemService{repo}
}

func (s *CreateListItemService) CreateListItem(listID int32, title domain.ItemTitle, description string, userID int32) (*domain.ListItem, error) {
	foundList, err := s.repo.FindListByID(listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting the user list", InternalError: err}
	}

	if foundList == nil {
		return nil, &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	maxPosition := int32(-1)

	if foundList.ItemsCount > 0 {
		maxPosition, err = s.repo.GetListItemsMaxPosition(listID, userID)
		if err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error getting the max position", InternalError: err}
		}
	}

	item := domain.ListItem{
		Title:       title,
		Description: description,
		ListID:      listID,
		UserID:      userID,
		Position:    maxPosition + 1,
	}

	err = s.repo.CreateListItem(&item)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the list item", InternalError: err}
	}

	return &item, nil
}
