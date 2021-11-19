package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
)

type GetListService struct {
	repo domain.ListsRepository
}

func NewGetListService(repo domain.ListsRepository) *GetListService {
	return &GetListService{repo}
}

func (s *GetListService) GetList(listID int32, userID int32) (*domain.List, error) {
	return s.repo.FindListByID(listID, userID)
}
