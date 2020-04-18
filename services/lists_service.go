package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
)

// ListsService is the service for the list entity
type ListsService struct {
	db *gorm.DB
}

// NewListsService returns a new lists service
func NewListsService(db *gorm.DB) ListsService {
	return ListsService{db}
}

// AddUserList  adds a list
func (s *ListsService) AddUserList(userID int32, l *models.List) (int32, error) {
	l.UserID = userID
	if err := s.db.Save(&l).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list", InternalError: err}
	}

	return l.ID, nil
}

// RemoveUserList removes a list
func (s *ListsService) RemoveUserList(id int32, userID int32) error {
	if err := s.db.Where(models.List{ID: id, UserID: userID}).Delete(models.List{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user list", InternalError: err}
	}
	return nil
}

// UpdateUserList updates an existing list
func (s *ListsService) UpdateUserList(id int32, userID int32, l *models.List) error {
	l.ID = id
	l.UserID = userID

	if err := s.db.Save(&l).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating list", InternalError: err}
	}

	return nil
}

// GetSingleUserList returns a single list from its id
func (s *ListsService) GetSingleUserList(id int32, userID int32, l *dtos.GetSingleListResultDto) error {
	if err := s.db.Where(models.List{ID: id, UserID: userID}).Preload("ListItems").Find(&l).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user list", InternalError: err}
	}

	// return s.listsRepository().GetOne(l, bson.D{{"_id", id}, {"userId", userID}}, nil)
	return nil
}

// GetUserLists returns the lists for the given user
func (s *ListsService) GetUserLists(userID int32, r *[]dtos.GetListsResultDto) error {
	if err := s.db.Where(models.List{UserID: userID}).Select("id,name").Find(&r).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user lists", InternalError: err}
	}
	return nil
}
