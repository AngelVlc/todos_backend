package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type GetAllListsService struct {
	repo domain.ListsRepository
}

func NewGetAllListsService(repo domain.ListsRepository) *GetAllListsService {
	return &GetAllListsService{repo}
}

func (s *GetAllListsService) GetAllLists(userID int32) ([]domain.List, error) {
	foundLists, err := s.repo.GetAllLists(userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all user lists", InternalError: err}
	}

	return foundLists, nil
}
