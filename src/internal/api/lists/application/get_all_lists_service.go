package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetAllListsService struct {
	repo domain.ListsRepository
}

func NewGetAllListsService(repo domain.ListsRepository) *GetAllListsService {
	return &GetAllListsService{repo}
}

func (s *GetAllListsService) GetAllLists(ctx context.Context, userID int32) ([]*domain.ListEntity, error) {
	foundLists, err := s.repo.GetAllListsForUser(ctx, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all user lists", InternalError: err}
	}

	return foundLists, nil
}
