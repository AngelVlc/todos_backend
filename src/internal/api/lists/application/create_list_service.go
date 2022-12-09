package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CreateListService struct {
	repo domain.ListsRepository
}

func NewCreateListService(repo domain.ListsRepository) *CreateListService {
	return &CreateListService{repo}
}

func (s *CreateListService) CreateList(ctx context.Context, name domain.ListNameValueObject, userID int32) (*domain.ListEntity, error) {
	err := name.CheckIfAlreadyExists(ctx, userID, s.repo)
	if err != nil {
		return nil, err
	}

	list := domain.ListEntity{
		Name:   name,
		UserID: userID,
	}

	err = s.repo.CreateList(ctx, &list)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user list", InternalError: err}
	}

	return &list, nil
}
