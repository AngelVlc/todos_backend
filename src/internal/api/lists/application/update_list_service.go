package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
)

type UpdateListService struct {
	repo     domain.ListsRepository
	eventBus events.EventBus
}

func NewUpdateListService(repo domain.ListsRepository, eventBus events.EventBus) *UpdateListService {
	return &UpdateListService{repo, eventBus}
}

func (s *UpdateListService) UpdateList(ctx context.Context, listRecordToUpdate *domain.ListRecord) error {
	foundList, err := s.repo.FindList(ctx, &domain.ListRecord{ID: listRecordToUpdate.ID, UserID: listRecordToUpdate.UserID})
	if err != nil {
		return err
	}

	if foundList.Name != listRecordToUpdate.Name {
		if existsList, err := s.repo.ExistsList(ctx, &domain.ListRecord{Name: listRecordToUpdate.Name, UserID: listRecordToUpdate.UserID}); err != nil {
			return &appErrors.UnexpectedError{Msg: "Error checking if a list with the same name already exists", InternalError: err}
		} else if existsList {
			return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
		}
	}

	if err := s.repo.UpdateList(ctx, listRecordToUpdate); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating the user list", InternalError: err}
	}

	go s.eventBus.Publish("listCreatedOrUpdated", listRecordToUpdate.ID)

	return nil
}
