package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"gorm.io/gorm"
)

type MySqlListsRepository struct {
	db *gorm.DB
}

func NewMySqlListsRepository(db *gorm.DB) *MySqlListsRepository {
	return &MySqlListsRepository{db}
}

func (r *MySqlListsRepository) FindList(ctx context.Context, query *domain.ListRecord) (*domain.ListRecord, error) {
	found := domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(query).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Take(&found).Error; err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) ExistsList(ctx context.Context, query *domain.ListRecord) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.ListRecord{}).Where(query).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) GetAllLists(ctx context.Context, userID int32) ([]domain.ListRecord, error) {
	found := []domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(domain.ListRecord{UserID: userID}).Find(&found).Error; err != nil {
		return nil, err
	}
	return found, nil
}

func (r *MySqlListsRepository) CreateList(ctx context.Context, list *domain.ListRecord) error {
	return r.db.WithContext(ctx).Create(list).Error
}

func (r *MySqlListsRepository) DeleteList(ctx context.Context, list *domain.ListRecord) error {
	return r.db.WithContext(ctx).Select("Items").Delete(list).Error
}

func (r *MySqlListsRepository) UpdateList(ctx context.Context, list *domain.ListRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(list).Error; err != nil {
			return err
		}

		var currentItems []int32
		for _, v := range list.Items {
			currentItems = append(currentItems, v.ID)
		}

		if err := tx.WithContext(ctx).Not(currentItems).Delete(&domain.ListItemRecord{}, "listId = ?", list.ID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *MySqlListsRepository) UpdateListItemsCounter(ctx context.Context, listID int32) error {
	subquery := r.db.WithContext(ctx).Model(&domain.ListItemRecord{}).Where(&domain.ListItemRecord{ListID: listID}).Select("COUNT(id)")

	return r.db.WithContext(ctx).Model(domain.ListRecord{}).Where(domain.ListRecord{ID: listID}).UpdateColumn("itemsCount", subquery).Error
}
