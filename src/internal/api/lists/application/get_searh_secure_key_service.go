package application

import (
	"fmt"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
)

type GetSearchSecureKeyService struct {
	searchClient search.SearchIndexClient
}

func NewGetSearchSecureKeyService(searchClient search.SearchIndexClient) *GetSearchSecureKeyService {
	return &GetSearchSecureKeyService{searchClient}
}

func (s *GetSearchSecureKeyService) GetSearchSecureKeyService(userID int32) (string, error) {
	filter := fmt.Sprintf("userID:%v", userID)
	k, err := s.searchClient.GenerateSecuredApiKey(filter)
	if err != nil {
		return "", &errors.UnexpectedError{Msg: "Error getting the search key", InternalError: err}
	}

	return k, nil
}
