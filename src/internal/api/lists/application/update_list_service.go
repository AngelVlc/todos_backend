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

func NewUpdateListService(listRepo domain.ListsRepository, eventBus events.EventBus) *UpdateListService {
	return &UpdateListService{listRepo, eventBus}
}

func (s *UpdateListService) UpdateList(ctx context.Context, listToUpdate *domain.ListEntity) (*domain.ListEntity, error) {
	foundList, err := s.repo.FindList(ctx, domain.ListEntity{ID: listToUpdate.ID, UserID: listToUpdate.UserID})
	if err != nil {
		return nil, err
	}

	if foundList.Name != listToUpdate.Name {
		if existsList, err := s.repo.ExistsList(ctx, domain.ListEntity{Name: listToUpdate.Name, UserID: listToUpdate.UserID}); err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error checking if a list with the same name already exists", InternalError: err}
		} else if existsList {
			return nil, &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
		}
	}

	updatedList, err := s.repo.UpdateList(ctx, listToUpdate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user list", InternalError: err}
	}

	go s.eventBus.Publish(events.ListUpdated, listToUpdate.ID)

	return updatedList, nil
}
