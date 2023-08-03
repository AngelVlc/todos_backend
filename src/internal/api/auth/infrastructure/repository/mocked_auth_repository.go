package repository

import (
	"context"
	"sync"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	"github.com/stretchr/testify/mock"
)

type MockedAuthRepository struct {
	mock.Mock
	mu sync.Mutex
	Wg sync.WaitGroup
}

func NewMockedAuthRepository() *MockedAuthRepository {
	return &MockedAuthRepository{}
}

func (m *MockedAuthRepository) ExistsRefreshToken(ctx context.Context, query domain.RefreshTokenEntity) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedAuthRepository) CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *domain.RefreshTokenEntity) error {
	defer m.Wg.Done()

	m.mu.Lock()
	defer m.mu.Unlock()

	args := m.Called(ctx, refreshToken)

	return args.Error(0)
}

func (m *MockedAuthRepository) DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error {
	args := m.Called(ctx, expTime)

	return args.Error(0)
}

func (m *MockedAuthRepository) GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]*domain.RefreshTokenEntity, error) {
	args := m.Called(ctx, paginationInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.RefreshTokenEntity), args.Error(1)
}

func (m *MockedAuthRepository) DeleteRefreshTokensByID(ctx context.Context, ids []int32) error {
	args := m.Called(ctx, ids)

	return args.Error(0)
}
