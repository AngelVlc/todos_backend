package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type GetListService struct {
	repo domain.ListsRepository
}

func NewGetListService(repo domain.ListsRepository) *GetListService {
	return &GetListService{repo}
}

func (s *GetListService) GetList(ctx context.Context, listID int32, userID int32) (*domain.List, error) {
	return s.repo.FindListByID(ctx, listID, userID)
}
