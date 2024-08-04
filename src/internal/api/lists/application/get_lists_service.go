package application

import (
	"context"
	"database/sql"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetListsService struct {
	repo domain.ListsRepository
}

func NewGetListsService(repo domain.ListsRepository) *GetListsService {
	return &GetListsService{repo}
}

func (s *GetListsService) GetLists(ctx context.Context, userID int32, categoryId *int32) ([]*domain.ListEntity, error) {
	query := domain.ListRecord{UserID: userID}
	if categoryId != nil {
		query.CategoryID = &sql.NullInt32{Int32: *categoryId, Valid: true}
	}

	foundLists, err := s.repo.GetLists(ctx, query)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all user lists", InternalError: err}
	}

	return foundLists.ToListEntities(), nil
}
