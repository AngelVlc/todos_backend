package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"gorm.io/gorm"
)

type MySqlCategoriesRepository struct {
	db *gorm.DB
}

func NewMySqlCategoriesRepository(db *gorm.DB) *MySqlCategoriesRepository {
	return &MySqlCategoriesRepository{db}
}

func (r *MySqlCategoriesRepository) FindCategory(ctx context.Context, query domain.CategoryEntity) (*domain.CategoryEntity, error) {
	foundCategory := domain.CategoryRecord{}
	if err := r.db.WithContext(ctx).Where(query.ToCategoryRecord()).Take(&foundCategory).Error; err != nil {
		return nil, err
	}

	return foundCategory.ToCategoryEntity(), nil
}

func (r *MySqlCategoriesRepository) ExistsCategory(ctx context.Context, query domain.CategoryEntity) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.CategoryRecord{}).Where(query.ToCategoryRecord()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlCategoriesRepository) GetAllCategoriesForUser(ctx context.Context, userID int32) ([]*domain.CategoryEntity, error) {
	foundCategories := []domain.CategoryRecord{}

	if err := r.db.WithContext(ctx).Where(domain.CategoryRecord{UserID: userID}).Find(&foundCategories).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.CategoryEntity, len(foundCategories))

	for i, l := range foundCategories {
		res[i] = l.ToCategoryEntity()
	}

	return res, nil
}

func (r *MySqlCategoriesRepository) CreateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	record := list.ToCategoryRecord()

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, err
	}

	return record.ToCategoryEntity(), nil
}

func (r *MySqlCategoriesRepository) DeleteCategory(ctx context.Context, query domain.CategoryEntity) error {
	return r.db.WithContext(ctx).Delete(query.ToCategoryRecord()).Error
}

func (r *MySqlCategoriesRepository) UpdateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	record := list.ToCategoryRecord()

	if err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error; err != nil {
		return nil, err
	}

	return record.ToCategoryEntity(), nil
}
