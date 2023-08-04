package repository

import (
	"context"
	"sync"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySqlAuthRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewMySqlAuthRepository(db *gorm.DB) *MySqlAuthRepository {
	return &MySqlAuthRepository{db, sync.Mutex{}}
}

func (r *MySqlAuthRepository) ExistsRefreshToken(ctx context.Context, query domain.RefreshTokenEntity) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.RefreshTokenRecord{}).Where(query.ToRefreshTokenRecord()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlAuthRepository) CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *domain.RefreshTokenEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "userId"}, {Name: "refreshToken"}},
		DoNothing: true,
	}).Create(refreshToken.ToRefreshTokenRecord()).Error
}

func (r *MySqlAuthRepository) DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error {
	return r.db.WithContext(ctx).Delete(domain.RefreshTokenRecord{}, "expirationDate <= ?", expTime).Error
}

func (r *MySqlAuthRepository) GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]*domain.RefreshTokenEntity, error) {
	foundRts := []domain.RefreshTokenRecord{}
	if err := r.db.WithContext(ctx).
		Select("id,userId,expirationDate").
		Limit(paginationInfo.Limit).Offset(paginationInfo.Offset).
		Order(paginationInfo.Order).
		Find(&foundRts).
		Error; err != nil {
		return nil, err
	}

	res := make([]*domain.RefreshTokenEntity, len(foundRts))

	for i, u := range foundRts {
		res[i] = u.ToRefreshTokenEntity()
	}

	return res, nil
}

func (r *MySqlAuthRepository) DeleteRefreshTokensByID(ctx context.Context, ids []int32) error {
	if err := r.db.WithContext(ctx).Delete(domain.RefreshTokenRecord{}, ids).Error; err != nil {
		return err
	}

	return nil
}
