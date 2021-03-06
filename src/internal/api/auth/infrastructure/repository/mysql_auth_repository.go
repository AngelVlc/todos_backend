package repository

import (
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/jinzhu/gorm"
)

type MySqlAuthRepository struct {
	db *gorm.DB
}

func NewMySqlAuthRepository(db *gorm.DB) *MySqlAuthRepository {
	return &MySqlAuthRepository{db}
}

func (r *MySqlAuthRepository) FindUserByName(userName domain.UserName) (*domain.User, error) {
	foundUser := domain.User{}
	err := r.db.Where(domain.User{Name: userName}).First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlAuthRepository) FindUserByID(userID int32) (*domain.User, error) {
	foundUser := domain.User{}
	err := r.db.Where(domain.User{ID: userID}).First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *MySqlAuthRepository) GetAllUsers() ([]domain.User, error) {
	res := []domain.User{}
	if err := r.db.Select("id,name,is_admin").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlAuthRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *MySqlAuthRepository) DeleteUser(userID int32) error {
	return r.db.Delete(domain.User{ID: userID}).Error
}

func (r *MySqlAuthRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(&user).Error
}

func (r *MySqlAuthRepository) FindRefreshTokenForUser(refreshToken string, userID int32) (*domain.RefreshToken, error) {
	found := domain.RefreshToken{}
	err := r.db.Where(domain.RefreshToken{RefreshToken: refreshToken, UserID: userID}).First(&found).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlAuthRepository) CreateRefreshToken(refreshToken *domain.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *MySqlAuthRepository) DeleteExpiredRefreshTokens(expTime time.Time) error {
	if err := r.db.Delete(domain.RefreshToken{}, "expirationDate <= ?", expTime).Error; err != nil {
		return err
	}
	return nil
}

func (r *MySqlAuthRepository) GetAllRefreshTokens() ([]domain.RefreshToken, error) {
	res := []domain.RefreshToken{}
	if err := r.db.Select("id,userId,expirationDate").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlAuthRepository) DeleteRefreshTokensByID(ids []int32) error {
	if err := r.db.Delete(domain.RefreshToken{}, ids).Error; err != nil {
		return err
	}
	return nil
}
