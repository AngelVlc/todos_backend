package repository

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySqlListsRepository struct {
	db *gorm.DB
}

func NewMySqlListsRepository(db *gorm.DB) *MySqlListsRepository {
	return &MySqlListsRepository{db}
}

func (r *MySqlListsRepository) ExistsList(name domain.ListName, userID int32) (bool, error) {
	count := int64(0)
	err := r.db.Model(&domain.List{}).Where(domain.List{Name: name, UserID: userID}).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MySqlListsRepository) FindListByID(listID int32, userID int32) (*domain.List, error) {
	found := domain.List{}
	err := r.db.Where(domain.List{ID: listID, UserID: userID}).First(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) GetAllLists(userID int32) ([]domain.List, error) {
	res := []domain.List{}
	if err := r.db.Where(domain.List{UserID: userID}).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlListsRepository) CreateList(list *domain.List) error {
	return r.db.Create(list).Error
}

func (r *MySqlListsRepository) DeleteList(listID int32, userID int32) error {
	return r.db.Where(domain.List{ID: listID, UserID: userID}).Delete(domain.List{}).Error
}

func (r *MySqlListsRepository) UpdateList(list *domain.List) error {
	return r.db.Model(list).Updates(domain.List{Name: list.Name}).Error
}

func (r *MySqlListsRepository) IncrementListCounter(listID int32) error {
	return r.db.Model(domain.List{}).Where(domain.List{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount + ?", 1)).Error
}

func (r *MySqlListsRepository) DecrementListCounter(listID int32) error {
	return r.db.Model(domain.List{}).Where(domain.List{ID: listID}).UpdateColumn("itemsCount", gorm.Expr("itemsCount - ?", 1)).Error
}

func (r *MySqlListsRepository) FindListItemByID(itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	found := domain.ListItem{}
	err := r.db.Where(domain.ListItem{ID: itemID, ListID: listID, UserID: userID}).First(&found).Error

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) GetAllListItems(listID int32, userID int32) ([]domain.ListItem, error) {
	res := []domain.ListItem{}
	if err := r.db.Where(domain.ListItem{ListID: listID, UserID: userID}).Order("position").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MySqlListsRepository) CreateListItem(listItem *domain.ListItem) error {
	return r.db.Create(listItem).Error
}

func (r *MySqlListsRepository) DeleteListItem(itemID int32, listID int32, userID int32) error {
	return r.db.Where(domain.ListItem{ID: itemID, ListID: listID, UserID: userID}).Delete(domain.ListItem{}).Error
}

func (r *MySqlListsRepository) UpdateListItem(listItem *domain.ListItem) error {
	return r.db.Save(&listItem).Error
}

func (r *MySqlListsRepository) BulkUpdateListItems(listItems []domain.ListItem) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"position"}),
	}).Create(listItems).Error
}

func (r *MySqlListsRepository) GetListItemsMaxPosition(listID int32, userID int32) (int32, error) {
	res := int32(-1)
	if err := r.db.Table("listItems").Where(domain.ListItem{ListID: listID, UserID: userID}).Select("MAX(position)").Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
