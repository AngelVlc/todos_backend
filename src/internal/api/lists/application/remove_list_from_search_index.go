package application

import (
	"context"
	"fmt"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
)

type RemoveListFromSearchIndexService struct {
	searchClient search.SearchIndexClient
}

func NewRemoveListFromSearchIndexService(searchClient search.SearchIndexClient) *RemoveListFromSearchIndexService {
	return &RemoveListFromSearchIndexService{searchClient}
}

func (s *RemoveListFromSearchIndexService) RemoveListFromSearchIndexService(ctx context.Context, listID int32) error {
	return s.searchClient.DeleteObject(fmt.Sprint(listID))
}
