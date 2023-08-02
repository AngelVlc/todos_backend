package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
)

type CreateListService struct {
	repo     domain.ListsRepository
	eventBus events.EventBus
}

func NewCreateListService(repo domain.ListsRepository, eventBus events.EventBus) *CreateListService {
	return &CreateListService{repo, eventBus}
}

func (s *CreateListService) CreateList(ctx context.Context, listRecordToCreate *domain.ListRecord) error {
	if existsList, err := s.repo.ExistsList(ctx, &domain.ListRecord{Name: listRecordToCreate.Name, UserID: listRecordToCreate.UserID}); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error checking if a list with the same name already exists", InternalError: err}
	} else if existsList {
		return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
	}

	err := s.repo.CreateList(ctx, listRecordToCreate)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error creating the user list", InternalError: err}
	}

	go s.eventBus.Publish("listCreatedOrUpdated", listRecordToCreate.ID)

	return nil
}
