package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/internal/api/config/domain"
	"gorm.io/gorm"
)

type MySqlConfigRepository struct {
	db *gorm.DB
}

func NewMySqlConfigRepository(db *gorm.DB) *MySqlConfigRepository {
	return &MySqlConfigRepository{db}
}

func (r *MySqlConfigRepository) ExistsAllowedOrigin(ctx context.Context, origin domain.Origin) (bool, error) {
	count := int64(0)
	err := r.db.WithContext(ctx).Model(&domain.AllowedOrigin{}).Where(domain.AllowedOrigin{Origin: origin}).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlConfigRepository) GetAllAllowedOrigins(ctx context.Context) ([]domain.AllowedOrigin, error) {
	res := []domain.AllowedOrigin{}
	if err := r.db.WithContext(ctx).Select("id,origin").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlConfigRepository) CreateAllowedOrigin(ctx context.Context, allowedOrigin *domain.AllowedOrigin) error {
	return r.db.WithContext(ctx).Create(allowedOrigin).Error
}

func (r *MySqlConfigRepository) DeleteAllowedOrigin(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(domain.AllowedOrigin{ID: id}).Error
}
