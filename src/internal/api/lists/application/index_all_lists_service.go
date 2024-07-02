package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
)

type IndexAllListsService struct {
	repo         domain.ListsRepository
	searchClient search.SearchIndexClient
}

func NewIndexAllListsService(repo domain.ListsRepository, searchClient search.SearchIndexClient) *IndexAllListsService {
	return &IndexAllListsService{repo, searchClient}
}

func (s *IndexAllListsService) IndexAllLists(ctx context.Context) error {
	foundLists, err := s.repo.GetLists(ctx, domain.ListRecord{})
	if err != nil {
		return &errors.UnexpectedError{Msg: "Error getting the lists", InternalError: err}
	}

	documents := make([]domain.ListSearchDocument, len(foundLists))

	for i, l := range foundLists {
		listEntity := l.ToListEntity()
		documents[i] = listEntity.ToListSearchDocument()
	}

	return s.searchClient.SaveObjects(documents)
}
