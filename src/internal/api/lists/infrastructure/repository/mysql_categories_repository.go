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

func (r *MySqlCategoriesRepository) FindCategory(ctx context.Context, query domain.CategoryRecord) (*domain.CategoryRecord, error) {
	foundCategory := domain.CategoryRecord{}
	if err := r.db.WithContext(ctx).Where(query).Take(&foundCategory).Error; err != nil {
		return nil, err
	}

	return &foundCategory, nil
}

func (r *MySqlCategoriesRepository) ExistsCategory(ctx context.Context, query domain.CategoryRecord) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.CategoryRecord{}).Where(query).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlCategoriesRepository) GetCategories(ctx context.Context, query domain.CategoryRecord) (domain.CategoryRecords, error) {
	foundCategories := []domain.CategoryRecord{}

	if err := r.db.WithContext(ctx).Where(query).Find(&foundCategories).Error; err != nil {
		return nil, err
	}

	return foundCategories, nil
}

func (r *MySqlCategoriesRepository) CreateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	record := list.ToCategoryRecord()

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, err
	}

	return record.ToCategoryEntity(), nil
}

func (r *MySqlCategoriesRepository) DeleteCategory(ctx context.Context, query domain.CategoryRecord) error {
	return r.db.WithContext(ctx).Delete(query).Error
}

func (r *MySqlCategoriesRepository) UpdateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	record := list.ToCategoryRecord()

	if err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error; err != nil {
		return nil, err
	}

	return record.ToCategoryEntity(), nil
}
