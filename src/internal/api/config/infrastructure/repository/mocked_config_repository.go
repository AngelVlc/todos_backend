package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/internal/api/config/domain"
	"github.com/stretchr/testify/mock"
)

type MockedConfigRepository struct {
	mock.Mock
}

func NewMockedConfigRepository() *MockedConfigRepository {
	return &MockedConfigRepository{}
}

func (m *MockedConfigRepository) ExistsAllowedOrigin(ctx context.Context, origin domain.Origin) (bool, error) {
	args := m.Called(ctx, origin)

	return args.Bool(0), args.Error(1)
}

func (m *MockedConfigRepository) GetAllAllowedOrigins(ctx context.Context) ([]domain.AllowedOrigin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.AllowedOrigin), args.Error(1)
}

func (m *MockedConfigRepository) CreateAllowedOrigin(ctx context.Context, allowedOrigin *domain.AllowedOrigin) error {
	args := m.Called(ctx, allowedOrigin)

	return args.Error(0)
}

func (m *MockedConfigRepository) DeleteAllowedOrigin(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
