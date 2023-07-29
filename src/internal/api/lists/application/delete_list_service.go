package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type DeleteListService struct {
	repo domain.ListsRepository
}

func NewDeleteListService(repo domain.ListsRepository) *DeleteListService {
	return &DeleteListService{repo}
}

func (s *DeleteListService) DeleteList(ctx context.Context, listID int32, userID int32) error {
	_, err := s.repo.FindList(ctx, &domain.ListRecord{ID: listID, UserID: userID})
	if err != nil {
		return err
	}

	err = s.repo.DeleteList(ctx, listID, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user list", InternalError: err}
	}

	return nil
}
