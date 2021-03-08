package repository

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/jinzhu/gorm"
)

type MySqlAuthRepository struct {
	db *gorm.DB
}

func NewMySqlAuthRepository(db *gorm.DB) *MySqlAuthRepository {
	return &MySqlAuthRepository{db}
}

func (r *MySqlAuthRepository) FindUserByName(userName *domain.AuthUserName) (*domain.AuthUser, error) {
	foundUser := domain.AuthUser{}
	err := r.db.Where(domain.AuthUser{Name: *userName}).First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlAuthRepository) FindUserByID(userID *int32) (*domain.AuthUser, error) {
	foundUser := domain.AuthUser{}
	err := r.db.Where(domain.AuthUser{ID: *userID}).First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}
