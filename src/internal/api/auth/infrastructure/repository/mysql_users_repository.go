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

func (r *MySqlUsersRepository) FindUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := domain.User{}
	err := r.db.WithContext(ctx).Where(user).Take(&foundUser).Error

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlUsersRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	res := []domain.User{}
	if err := r.db.WithContext(ctx).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
