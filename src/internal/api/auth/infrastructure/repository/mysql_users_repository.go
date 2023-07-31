package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"gorm.io/gorm"
)

type MySqlUsersRepository struct {
	db *gorm.DB
}

func NewMySqlUsersRepository(db *gorm.DB) *MySqlUsersRepository {
	return &MySqlUsersRepository{db}
}

func (r *MySqlUsersRepository) FindUser(ctx context.Context, query *domain.UserRecord) (*domain.UserRecord, error) {
	foundUser := domain.UserRecord{}
	if err := r.db.WithContext(ctx).Where(query).Take(&foundUser).Error; err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlUsersRepository) ExistsUser(ctx context.Context, query *domain.UserRecord) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.UserRecord{}).Where(query).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlUsersRepository) GetAll(ctx context.Context) ([]domain.UserRecord, error) {
	res := []domain.UserRecord{}
	if err := r.db.WithContext(ctx).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MySqlUsersRepository) Create(ctx context.Context, user *domain.UserRecord) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *MySqlUsersRepository) Delete(ctx context.Context, query *domain.UserRecord) error {
	return r.db.WithContext(ctx).Delete(query).Error
}

func (r *MySqlUsersRepository) Update(ctx context.Context, user *domain.UserRecord) error {
	return r.db.WithContext(ctx).Save(&user).Error
}
