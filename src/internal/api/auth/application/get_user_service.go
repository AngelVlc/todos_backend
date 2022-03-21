package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/internal/api/auth/domain"
)

type GetUserService struct {
	repo domain.AuthRepository
}

func NewGetUserService(repo domain.AuthRepository) *GetUserService {
	return &GetUserService{repo}
}

func (s *GetUserService) GetUser(ctx context.Context, userID int32) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, userID)
}
