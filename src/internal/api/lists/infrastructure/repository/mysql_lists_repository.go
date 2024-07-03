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

func orderItems(db *gorm.DB) *gorm.DB {
	return db.Order("position ASC")
}

func (r *MySqlListsRepository) FindList(ctx context.Context, query domain.ListRecord) (*domain.ListRecord, error) {
	foundList := domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(query).Preload("Items", orderItems).Take(&foundList).Error; err != nil {
		return nil, err
	}

	return &foundList, nil
}

func (r *MySqlListsRepository) ExistsList(ctx context.Context, query domain.ListRecord) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.ListRecord{}).Where(query).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) GetLists(ctx context.Context, query domain.ListRecord) (domain.ListRecords, error) {
	foundLists := []domain.ListRecord{}

	if err := r.db.WithContext(ctx).Where(query).Find(&foundLists).Error; err != nil {
		return nil, err
	}

	return foundLists, nil
}

func (r *MySqlListsRepository) CreateList(ctx context.Context, record *domain.ListRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	return nil
}

func (r *MySqlListsRepository) DeleteList(ctx context.Context, query domain.ListRecord) error {
	return r.db.WithContext(ctx).Select("Items").Delete(query).Error
}

func (r *MySqlListsRepository) UpdateList(ctx context.Context, record *domain.ListRecord) error {
	error := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(record).Error; err != nil {
			return err
		}

		var currentItems []int32
		for _, v := range record.Items {
			currentItems = append(currentItems, v.ID)
		}

		if err := tx.WithContext(ctx).Not(currentItems).Delete(&domain.ListItemRecord{}, "listId = ?", record.ID).Error; err != nil {
			return err
		}

		return nil
	})

	return error
}

func (r *MySqlListsRepository) UpdateListItemsCount(ctx context.Context, listID int32) error {
	subquery := r.db.WithContext(ctx).Model(&domain.ListItemRecord{}).Where(&domain.ListItemRecord{ListID: listID}).Select("COUNT(id)")

	return r.db.WithContext(ctx).Model(domain.ListRecord{}).Where(domain.ListRecord{ID: listID}).UpdateColumn("itemsCount", subquery).Error
}
