package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
)

type GetUserService struct {
	repo domain.AuthRepository
}

func NewGetUserService(repo domain.AuthRepository) *GetUserService {
	return &GetUserService{repo}
}

func (s *GetUserService) GetUser(userID int32) (*domain.User, error) {
	return s.repo.FindUserByID(userID)
}
