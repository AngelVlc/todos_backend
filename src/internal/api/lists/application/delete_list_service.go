package application

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type DeleteListService struct {
	repo domain.ListsRepository
}

func NewDeleteListService(repo domain.ListsRepository) *DeleteListService {
	return &DeleteListService{repo}
}

func (s *DeleteListService) DeleteList(ctx context.Context, listID int32, userID int32) error {
	_, err := s.repo.FindListByID(ctx, listID, userID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteList(listID, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user list", InternalError: err}
	}

	return nil
}
