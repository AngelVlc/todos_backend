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

func (r *MySqlUsersRepository) FindUser(ctx context.Context, filter *domain.UserEntity) (*domain.UserEntity, error) {
	foundUser := domain.UserEntity{}
	err := r.db.WithContext(ctx).Where(filter).Take(&foundUser).Error

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlUsersRepository) ExistsUser(ctx context.Context, filter *domain.UserEntity) (bool, error) {
	count := int64(0)
	err := r.db.WithContext(ctx).Model(&domain.UserEntity{}).Where(filter).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlUsersRepository) GetAll(ctx context.Context) ([]domain.UserEntity, error) {
	res := []domain.UserEntity{}
	if err := r.db.WithContext(ctx).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlUsersRepository) Create(ctx context.Context, user *domain.UserEntity) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *MySqlUsersRepository) Delete(ctx context.Context, filter *domain.UserEntity) error {
	return r.db.WithContext(ctx).Delete(filter).Error
}

func (r *MySqlUsersRepository) Update(ctx context.Context, user *domain.UserEntity) error {
	return r.db.WithContext(ctx).Save(&user).Error
}
