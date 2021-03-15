package repository

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/jinzhu/gorm"
)

type MySqlListsRepository struct {
	db *gorm.DB
}

func NewMySqlListsRepository(db *gorm.DB) *MySqlListsRepository {
	return &MySqlListsRepository{db}
}

func (r *MySqlListsRepository) FindListByID(listID int32, userID int32) (*domain.List, error) {
	found := domain.List{}
	err := r.db.Where(domain.List{ID: listID, UserID: userID}).First(&found).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

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
	return r.db.Save(&list).Error
}

func (r *MySqlListsRepository) FindListItemByID(itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	found := domain.ListItem{}
	err := r.db.Where(domain.ListItem{ID: itemID, ListID: listID, UserID: userID}).First(&found).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *MySqlListsRepository) GetAllItems(listID int32, userID int32) ([]domain.ListItem, error) {
	res := []domain.ListItem{}
	if err := r.db.Where(domain.ListItem{ListID: listID, UserID: userID}).Find(&res).Error; err != nil {
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
