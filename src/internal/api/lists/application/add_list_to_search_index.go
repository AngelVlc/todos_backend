package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
)

type AddListToSearchIndexService struct {
	repo         domain.ListsRepository
	searchClient search.SearchIndexClient
}

func NewAddListToSearchIndexService(repo domain.ListsRepository, searchClient search.SearchIndexClient) *AddListToSearchIndexService {
	return &AddListToSearchIndexService{repo, searchClient}
}

func (s *AddListToSearchIndexService) AddListToSearchIndexService(ctx context.Context, listID int32) error {
	foundList, err := s.repo.FindList(ctx, domain.ListEntity{ID: listID})
	if err != nil {
		return &errors.UnexpectedError{Msg: "Error getting the list", InternalError: err}
	}

	return s.searchClient.SaveObjects(foundList.ToListSearchDocument())
}
