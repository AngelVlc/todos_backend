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

func (r *MySqlUsersRepository) FindUser(ctx context.Context, query domain.UserRecord) (*domain.UserRecord, error) {
	foundUser := domain.UserRecord{}
	if err := r.db.WithContext(ctx).Where(query).Take(&foundUser).Error; err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlUsersRepository) ExistsUser(ctx context.Context, query domain.UserRecord) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.UserRecord{}).Where(query).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlUsersRepository) GetAll(ctx context.Context) ([]*domain.UserEntity, error) {
	foundUsers := []domain.UserRecord{}
	if err := r.db.WithContext(ctx).Find(&foundUsers).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.UserEntity, len(foundUsers))

	for i, u := range foundUsers {
		res[i] = u.ToUserEntity()
	}

	return res, nil
}

func (r *MySqlUsersRepository) Create(ctx context.Context, record *domain.UserRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *MySqlUsersRepository) Delete(ctx context.Context, query domain.UserRecord) error {
	return r.db.WithContext(ctx).Delete(query).Error
}

func (r *MySqlUsersRepository) Update(ctx context.Context, record *domain.UserRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}
