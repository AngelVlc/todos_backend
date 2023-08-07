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

func (r *MySqlListsRepository) FindList(ctx context.Context, query domain.ListEntity) (*domain.ListEntity, error) {
	foundList := domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(query.ToListRecord()).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Take(&foundList).Error; err != nil {
		return nil, err
	}

	return foundList.ToListEntity(), nil
}

func (r *MySqlListsRepository) ExistsList(ctx context.Context, query domain.ListEntity) (bool, error) {
	count := int64(0)
	if err := r.db.WithContext(ctx).Model(&domain.ListRecord{}).Where(query.ToListRecord()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) GetAllLists(ctx context.Context, userID int32) ([]*domain.ListEntity, error) {
	foundLists := []domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(domain.ListRecord{UserID: userID}).Find(&foundLists).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.ListEntity, len(foundLists))

	for i, l := range foundLists {
		res[i] = l.ToListEntity()
	}

	return res, nil
}

func (r *MySqlListsRepository) CreateList(ctx context.Context, list *domain.ListEntity) (*domain.ListEntity, error) {
	record := list.ToListRecord()

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return nil, err
	}

	return record.ToListEntity(), nil
}

func (r *MySqlListsRepository) DeleteList(ctx context.Context, query domain.ListEntity) error {
	return r.db.WithContext(ctx).Select("Items").Delete(query.ToListRecord()).Error
}

func (r *MySqlListsRepository) UpdateList(ctx context.Context, list *domain.ListEntity) (*domain.ListEntity, error) {
	record := list.ToListRecord()

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

	return record.ToListEntity(), error
}

func (r *MySqlListsRepository) UpdateListItemsCounter(ctx context.Context, listID int32) error {
	subquery := r.db.WithContext(ctx).Model(&domain.ListItemRecord{}).Where(&domain.ListItemRecord{ListID: listID}).Select("COUNT(id)")

	return r.db.WithContext(ctx).Model(domain.ListRecord{}).Where(domain.ListRecord{ID: listID}).UpdateColumn("itemsCount", subquery).Error
}
