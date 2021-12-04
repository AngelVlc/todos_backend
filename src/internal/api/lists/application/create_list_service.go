package application

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type CreateListService struct {
	repo domain.ListsRepository
}

func NewCreateListService(repo domain.ListsRepository) *CreateListService {
	return &CreateListService{repo}
}

func (s *CreateListService) CreateList(ctx context.Context, name domain.ListName, userID int32) (*domain.List, error) {
	err := name.CheckIfAlreadyExists(ctx, userID, s.repo)
	if err != nil {
		return nil, err
	}

	list := domain.List{
		Name:   name,
		UserID: userID,
	}

	err = s.repo.CreateList(ctx, &list)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user list", InternalError: err}
	}

	return &list, nil
}
