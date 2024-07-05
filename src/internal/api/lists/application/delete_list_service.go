package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
)

type DeleteListService struct {
	repo     domain.ListsRepository
	eventBus events.EventBus
}

func NewDeleteListService(repo domain.ListsRepository, eventBus events.EventBus) *DeleteListService {
	return &DeleteListService{repo, eventBus}
}

func (s *DeleteListService) DeleteList(ctx context.Context, listID int32, userID int32) error {
	foundList, err := s.repo.FindList(ctx, domain.ListRecord{ID: listID, UserID: userID})
	if err != nil {
		return err
	}

	if err := s.repo.DeleteList(ctx, *foundList); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user list", InternalError: err}
	}

	go s.eventBus.Publish(events.ListDeleted, listID)

	return nil
}
