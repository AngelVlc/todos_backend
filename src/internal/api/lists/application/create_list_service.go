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

func (s *CreateListService) CreateList(ctx context.Context, listToCreate *domain.ListEntity) error {
	if existsList, err := s.repo.ExistsList(ctx, domain.ListRecord{Name: listToCreate.Name.String(), UserID: listToCreate.UserID}); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error checking if a list with the same name already exists", InternalError: err}
	} else if existsList {
		return &appErrors.BadRequestError{Msg: "A list with the same name already exists", InternalError: nil}
	}

	record := listToCreate.ToListRecord()

	err := s.repo.CreateList(ctx, record)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error creating the user list", InternalError: err}
	}

	listToCreate.ID = record.ID

	go s.eventBus.Publish(events.ListCreated, record.ID)

	return nil
}
