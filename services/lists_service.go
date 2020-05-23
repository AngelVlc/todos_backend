package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsService interface {
	AddUserList(userID int32, l *models.List) (int32, error)
	RemoveUserList(id int32, userID int32) error
	UpdateUserList(id int32, userID int32, l *models.List) error
	GetSingleUserList(id int32, userID int32, l *dtos.GetSingleListResultDto) error
	GetUserLists(userID int32, r *[]dtos.GetListsResultDto) error
	GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error
	AddUserListItem(userId int32, i *models.ListItem) (int32, error)
	RemoveUserListItem(id int32, listID int32, userID int32) error
	UpdateUserListItem(id int32, listID int32, userID int32, i *models.ListItem) error
}

type MockedListsService struct {
	mock.Mock
}

func NewMockedListsService() *MockedListsService {
	return &MockedListsService{}
}

func (m *MockedListsService) AddUserList(userID int32, l *models.List) (int32, error) {
	args := m.Called(userID, l)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsService) RemoveUserList(id int32, userID int32) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockedListsService) UpdateUserList(id int32, userID int32, l *models.List) error {
	args := m.Called(id, userID, l)
	return args.Error(0)
}

func (m *MockedListsService) GetSingleUserList(id int32, userID int32, l *dtos.GetSingleListResultDto) error {
	args := m.Called(id, userID, l)
	return args.Error(0)
}

func (m *MockedListsService) GetUserLists(userID int32, r *[]dtos.GetListsResultDto) error {
	args := m.Called(userID, r)
	return args.Error(0)
}

func (m *MockedListsService) GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error {
	args := m.Called(id, listID, userID, i)
	return args.Error(0)
}

func (m *MockedListsService) AddUserListItem(userID int32, i *models.ListItem) (int32, error) {
	args := m.Called(userID, i)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsService) RemoveUserListItem(id int32, listID int32, userID int32) error {
	args := m.Called(id, listID, userID)
	return args.Error(0)
}

func (m *MockedListsService) UpdateUserListItem(id int32, listID int32, userID int32, i *models.ListItem) error {
	args := m.Called(id, listID, userID, i)
	return args.Error(0)
}

// DefaultListsService is the service for the list entity
type DefaultListsService struct {
	db *gorm.DB
}

// NewDefaultListsService returns a new lists service
func NewDefaultListsService(db *gorm.DB) *DefaultListsService {
	return &DefaultListsService{db}
}

// AddUserList  adds a list
func (s *DefaultListsService) AddUserList(userID int32, l *models.List) (int32, error) {
	l.UserID = userID
	if err := s.db.Create(&l).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list", InternalError: err}
	}

	return l.ID, nil
}

// RemoveUserList removes a list
func (s *DefaultListsService) RemoveUserList(id int32, userID int32) error {
	if err := s.db.Where(models.List{ID: id, UserID: userID}).Delete(models.List{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user list", InternalError: err}
	}
	return nil
}

// UpdateUserList updates an existing list
func (s *DefaultListsService) UpdateUserList(id int32, userID int32, l *models.List) error {
	l.ID = id
	l.UserID = userID

	if err := s.db.Save(&l).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating list", InternalError: err}
	}

	return nil
}

// GetSingleUserList returns a single list from its id
func (s *DefaultListsService) GetSingleUserList(id int32, userID int32, l *dtos.GetSingleListResultDto) error {
	if err := s.db.Where(models.List{ID: id, UserID: userID}).Preload("ListItems").Find(&l).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user list", InternalError: err}
	}

	return nil
}

// GetUserLists returns the lists for the given user
func (s *DefaultListsService) GetUserLists(userID int32, r *[]dtos.GetListsResultDto) error {
	if err := s.db.Where(models.List{UserID: userID}).Select("id,name").Find(&r).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user lists", InternalError: err}
	}
	return nil
}

// GetUserListItem returns a list item
func (s *DefaultListsService) GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error {
	if err := s.db.Joins("JOIN lists on listItems.listId=lists.id").Where(models.List{ID: listID, UserID: userID}).Where(models.ListItem{ID: id}).Find(&i).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user list item", InternalError: err}
	}

	return nil
}

// AddUserListItem adds a list item
func (s *DefaultListsService) AddUserListItem(userID int32, i *models.ListItem) (int32, error) {
	foundList := &dtos.GetSingleListResultDto{}

	if err := s.GetSingleUserList(i.ListID, userID, foundList); err != nil {
		return 0, err
	}

	if err := s.db.Create(&i).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list item", InternalError: err}
	}

	return i.ID, nil
}

// RemoveUserListItem removes a list item
func (s *DefaultListsService) RemoveUserListItem(id int32, listID int32, userID int32) error {
	foundList := &dtos.GetSingleListResultDto{}

	if err := s.GetSingleUserList(listID, userID, foundList); err != nil {
		return err
	}

	if err := s.db.Where(models.ListItem{ID: id, ListID: listID}).Delete(models.ListItem{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user list item", InternalError: err}
	}
	return nil
}

// UpdateUserListItem updates a list item
func (s *DefaultListsService) UpdateUserListItem(id int32, listID int32, userID int32, i *models.ListItem) error {
	foundList := &dtos.GetSingleListResultDto{}

	if err := s.GetSingleUserList(listID, userID, foundList); err != nil {
		return err
	}

	i.ID = id
	i.ListID = listID

	if err := s.db.Save(&i).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating list item", InternalError: err}
	}

	return nil
}
