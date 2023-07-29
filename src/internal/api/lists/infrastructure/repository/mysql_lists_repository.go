package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySqlListsRepository struct {
	db *gorm.DB
}

func NewMySqlListsRepository(db *gorm.DB) *MySqlListsRepository {
	return &MySqlListsRepository{db}
}

func (r *MySqlListsRepository) FindList(ctx context.Context, list *domain.ListRecord) (*domain.ListRecord, error) {
	found := domain.ListRecord{}
	err := r.db.WithContext(ctx).Where(list).Take(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) ExistsList(ctx context.Context, list *domain.ListRecord) (bool, error) {
	count := int64(0)
	err := r.db.WithContext(ctx).Model(&domain.ListRecord{}).Where(list).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) GetAllLists(ctx context.Context, userID int32) ([]domain.ListRecord, error) {
	res := []domain.ListRecord{}
	if err := r.db.WithContext(ctx).Where(domain.ListRecord{UserID: userID}).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlListsRepository) CreateList(ctx context.Context, list *domain.ListRecord) error {
	return r.db.WithContext(ctx).Create(list).Error
}

func (r *MySqlListsRepository) DeleteList(ctx context.Context, listID int32, userID int32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where(domain.ListItemRecord{ListID: listID, UserID: userID}).Delete(domain.ListItemRecord{}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Where(domain.ListRecord{ID: listID, UserID: userID}).Delete(domain.ListRecord{}).Error
	})
}

func (r *MySqlListsRepository) UpdateList(ctx context.Context, list *domain.ListRecord) error {
	return r.db.WithContext(ctx).Model(list).Updates(domain.ListRecord{Name: list.Name}).Error
}

func (r *MySqlListsRepository) IncrementListCounter(ctx context.Context, listID int32) error {
	return r.db.WithContext(ctx).Model(domain.ListRecord{}).Where(domain.ListRecord{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount + ?", 1)).Error
}

func (r *MySqlListsRepository) DecrementListCounter(ctx context.Context, listID int32) error {
	return r.db.WithContext(ctx).Model(domain.ListRecord{}).Where(domain.ListRecord{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount - ?", 1)).Error
}

func (r *MySqlListsRepository) FindListItem(ctx context.Context, listItem *domain.ListItemRecord) (*domain.ListItemRecord, error) {
	found := domain.ListItemRecord{}
	err := r.db.WithContext(ctx).Where(listItem).Take(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) GetAllListItems(ctx context.Context, listID int32, userID int32) ([]domain.ListItemRecord, error) {
	res := []domain.ListItemRecord{}
	if err := r.db.WithContext(ctx).Where(domain.ListItemRecord{ListID: listID, UserID: userID}).Order("position").Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MySqlListsRepository) CreateListItem(ctx context.Context, listItem *domain.ListItemRecord) error {
	return r.db.WithContext(ctx).Create(listItem).Error
}

func (r *MySqlListsRepository) DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error {
	return r.db.WithContext(ctx).Where(domain.ListItemRecord{ID: itemID, ListID: listID, UserID: userID}).Delete(domain.ListItemRecord{}).Error
}

func (r *MySqlListsRepository) UpdateListItem(ctx context.Context, listItem *domain.ListItemRecord) error {
	return r.db.WithContext(ctx).Save(&listItem).Error
}

func (r *MySqlListsRepository) BulkUpdateListItems(ctx context.Context, listItems []domain.ListItemRecord) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"position"}),
	}).Create(listItems).Error
}

func (r *MySqlListsRepository) GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error) {
	res := int32(-1)
	if err := r.db.WithContext(ctx).Table("listItems").Where(domain.ListItemRecord{ListID: listID, UserID: userID}).Select("MAX(position)").Scan(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
