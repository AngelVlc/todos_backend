package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
)

type GetUserService struct {
	repo domain.UsersRepository
}

func NewGetUserService(repo domain.UsersRepository) *GetUserService {
	return &GetUserService{repo}
}

func (s *GetUserService) GetUser(ctx context.Context, userID int32) (*domain.UserEntity, error) {
	return s.repo.FindUser(ctx, &domain.UserEntity{ID: userID})
}
