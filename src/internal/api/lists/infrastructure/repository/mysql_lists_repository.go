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

func (r *MySqlListsRepository) FindList(ctx context.Context, list *domain.List) (*domain.List, error) {
	found := domain.List{}
	err := r.db.WithContext(ctx).Where(list).Take(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) ExistsList(ctx context.Context, list *domain.List) (bool, error) {
	count := int64(0)
	err := r.db.WithContext(ctx).Model(&domain.List{}).Where(list).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) GetAllLists(ctx context.Context, userID int32) ([]domain.List, error) {
	res := []domain.List{}
	if err := r.db.WithContext(ctx).Where(domain.List{UserID: userID}).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlListsRepository) CreateList(ctx context.Context, list *domain.List) error {
	return r.db.WithContext(ctx).Select("Name", "UserID", "IsQuickList").Create(list).Error
}

func (r *MySqlListsRepository) DeleteList(ctx context.Context, listID int32, userID int32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where(domain.ListItem{ListID: listID, UserID: userID}).Delete(domain.ListItem{}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Where(domain.List{ID: listID, UserID: userID}).Delete(domain.List{}).Error
	})
}

func (r *MySqlListsRepository) UpdateList(ctx context.Context, list *domain.List) error {
	return r.db.WithContext(ctx).Model(list).Select("Name", "IsQuickList").Updates(domain.List{Name: list.Name, IsQuickList: list.IsQuickList}).Error
}

func (r *MySqlListsRepository) IncrementListCounter(ctx context.Context, listID int32) error {
	return r.db.WithContext(ctx).Model(domain.List{}).Where(domain.List{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount + ?", 1)).Error
}

func (r *MySqlListsRepository) DecrementListCounter(ctx context.Context, listID int32) error {
	return r.db.WithContext(ctx).Model(domain.List{}).Where(domain.List{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount - ?", 1)).Error
}

func (r *MySqlListsRepository) FindListItem(ctx context.Context, listItem *domain.ListItem) (*domain.ListItem, error) {
	found := domain.ListItem{}
	err := r.db.WithContext(ctx).Where(listItem).Take(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) GetAllListItems(ctx context.Context, listID int32, userID int32) ([]domain.ListItem, error) {
	res := []domain.ListItem{}
	if err := r.db.WithContext(ctx).Where(domain.ListItem{ListID: listID, UserID: userID}).Order("position").Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MySqlListsRepository) CreateListItem(ctx context.Context, listItem *domain.ListItem) error {
	return r.db.WithContext(ctx).Create(listItem).Error
}

func (r *MySqlListsRepository) DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error {
	return r.db.WithContext(ctx).Where(domain.ListItem{ID: itemID, ListID: listID, UserID: userID}).Delete(domain.ListItem{}).Error
}

func (r *MySqlListsRepository) UpdateListItem(ctx context.Context, listItem *domain.ListItem) error {
	return r.db.WithContext(ctx).Save(&listItem).Error
}

func (r *MySqlListsRepository) BulkUpdateListItems(ctx context.Context, listItems []domain.ListItem) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"position"}),
	}).Create(listItems).Error
}

func (r *MySqlListsRepository) GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error) {
	res := int32(-1)
	if err := r.db.WithContext(ctx).Table("listItems").Where(domain.ListItem{ListID: listID, UserID: userID}).Select("MAX(position)").Scan(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
