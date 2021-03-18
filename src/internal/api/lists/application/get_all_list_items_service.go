package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type GetAllListItemsService struct {
	repo domain.ListsRepository
}

func NewGetAllListItemsService(repo domain.ListsRepository) *GetAllListItemsService {
	return &GetAllListItemsService{repo}
}

func (s *GetAllListItemsService) GetAllListItems(listID int32, userID int32) ([]domain.ListItem, error) {
	foundLists, err := s.repo.GetAllListItems(listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all list items", InternalError: err}
	}

	return foundLists, nil
}
