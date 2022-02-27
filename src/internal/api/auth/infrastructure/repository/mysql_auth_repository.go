package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos/internal/api/shared/domain"
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

func (r *MySqlAuthRepository) ExistsUser(ctx context.Context, userName domain.UserName) (bool, error) {
	count := int64(0)
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where(domain.User{Name: userName}).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlAuthRepository) FindUserByName(ctx context.Context, userName domain.UserName) (*domain.User, error) {
	return r.findUser(ctx, domain.User{Name: userName})
}

func (r *MySqlAuthRepository) FindUserByID(ctx context.Context, userID int32) (*domain.User, error) {
	return r.findUser(ctx, domain.User{ID: userID})
}

func (r *MySqlAuthRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	res := []domain.User{}
	if err := r.db.WithContext(ctx).Select("id,name,isAdmin").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlAuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *MySqlAuthRepository) DeleteUser(ctx context.Context, userID int32) error {
	return r.db.WithContext(ctx).Delete(domain.User{ID: userID}).Error
}

func (r *MySqlAuthRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *MySqlAuthRepository) FindRefreshTokenForUser(ctx context.Context, refreshToken string, userID int32) (*domain.RefreshToken, error) {
	found := domain.RefreshToken{}
	err := r.db.WithContext(ctx).Where(domain.RefreshToken{RefreshToken: refreshToken, UserID: userID}).Take(&found).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlAuthRepository) CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "userId"}, {Name: "refreshToken"}},
		DoNothing: true,
	}).Create(refreshToken).Error
}

func (r *MySqlAuthRepository) DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error {
	return r.db.WithContext(ctx).Delete(domain.RefreshToken{}, "expirationDate <= ?", expTime).Error
}

func (r *MySqlAuthRepository) GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]domain.RefreshToken, error) {
	res := []domain.RefreshToken{}
	if err := r.db.WithContext(ctx).
		Select("id,userId,expirationDate").
		Limit(paginationInfo.Limit).Offset(paginationInfo.Offset).
		Order(paginationInfo.Order).
		Find(&res).
		Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MySqlAuthRepository) DeleteRefreshTokensByID(ctx context.Context, ids []int32) error {
	if err := r.db.WithContext(ctx).Delete(domain.RefreshToken{}, ids).Error; err != nil {
		return err
	}

	return nil
}

func (r *MySqlAuthRepository) findUser(ctx context.Context, where domain.User) (*domain.User, error) {
	foundUser := domain.User{}
	err := r.db.WithContext(ctx).Where(where).Take(&foundUser).Error

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}
