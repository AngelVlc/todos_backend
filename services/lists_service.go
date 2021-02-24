package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsService interface {
	AddUserList(userID int32, dto *dtos.ListDto) (int32, error)
	RemoveUserList(id int32, userID int32) error
	UpdateUserList(id int32, userID int32, dto *dtos.ListDto) error
	GetUserList(id int32, userID int32) (*dtos.ListResponseDto, error)
	GetUserLists(userID int32) ([]*dtos.ListResponseDto, error)

	GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error
	AddUserListItem(listID int32, userId int32, dto *dtos.ListItemDto) (int32, error)
	RemoveUserListItem(id int32, listID int32, userID int32) error
	UpdateUserListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error
}

type MockedListsService struct {
	mock.Mock
}

func NewMockedListsService() *MockedListsService {
	return &MockedListsService{}
}

func (m *MockedListsService) AddUserList(userID int32, dto *dtos.ListDto) (int32, error) {
	args := m.Called(userID, dto)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsService) RemoveUserList(id int32, userID int32) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockedListsService) UpdateUserList(id int32, userID int32, dto *dtos.ListDto) error {
	args := m.Called(id, userID, dto)
	return args.Error(0)
}

func (m *MockedListsService) GetUserList(id int32, userID int32) (*dtos.ListResponseDto, error) {
	args := m.Called(id, userID)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.ListResponseDto), args.Error(1)
}

func (m *MockedListsService) GetUserLists(userID int32) ([]*dtos.ListResponseDto, error) {
	args := m.Called(userID)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dtos.ListResponseDto), args.Error(1)
}

func (m *MockedListsService) GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error {
	args := m.Called(id, listID, userID, i)
	return args.Error(0)
}

func (m *MockedListsService) AddUserListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error) {
	args := m.Called(listID, userID, dto)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsService) RemoveUserListItem(id int32, listID int32, userID int32) error {
	args := m.Called(id, listID, userID)
	return args.Error(0)
}

func (m *MockedListsService) UpdateUserListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error {
	args := m.Called(id, listID, userID, dto)
	return args.Error(0)
}

// DefaultListsService is the service for the list entity
type DefaultListsService struct {
	db        *gorm.DB
	listsRepo repositories.ListsRepository
}

// NewDefaultListsService returns a new lists service
func NewDefaultListsService(db *gorm.DB, listsRepo repositories.ListsRepository) *DefaultListsService {
	return &DefaultListsService{db, listsRepo}
}

// AddUserList  adds a list
func (s *DefaultListsService) AddUserList(userID int32, dto *dtos.ListDto) (int32, error) {
	l := models.List{}
	l.FromDto(dto)
	l.UserID = userID

	return s.listsRepo.Insert(&l)
}

// RemoveUserList removes a list
func (s *DefaultListsService) RemoveUserList(id int32, userID int32) error {
	return s.listsRepo.Remove(id, userID)
}

// UpdateUserList updates an existing list
func (s *DefaultListsService) UpdateUserList(id int32, userID int32, dto *dtos.ListDto) error {
	foundList, err := s.listsRepo.FindByID(id, userID)
	if err != nil {
		return err
	}

	if foundList == nil {
		return &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	foundList.FromDto(dto)

	return s.listsRepo.Update(foundList)
}

func (s *DefaultListsService) GetUserList(id int32, userID int32) (*dtos.ListResponseDto, error) {
	foundList, err := s.listsRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	if foundList == nil {
		return nil, nil
	}

	return foundList.ToResponseDto(), nil
}

// GetUserLists returns the lists for the given user
func (s *DefaultListsService) GetUserLists(userID int32) ([]*dtos.ListResponseDto, error) {
	found, err := s.listsRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	res := make([]*dtos.ListResponseDto, len(found))

	for i, v := range found {
		res[i] = v.ToResponseDto()
	}

	return res, nil
}

// GetUserListItem returns a list item
func (s *DefaultListsService) GetUserListItem(id int32, listID int32, userID int32, i *dtos.GetItemResultDto) error {
	if err := s.db.Joins("JOIN lists on listItems.listId=lists.id").Where(models.List{ID: listID, UserID: userID}).Where(models.ListItem{ID: id}).Find(&i).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user list item", InternalError: err}
	}

	return nil
}

// TEMP
func (s *DefaultListsService) getListItem(id int32, listID int32, userID int32) (*models.ListItem, error) {
	foundItem := models.ListItem{}

	if err := s.db.Joins("JOIN lists on listItems.listId=lists.id").Where(models.List{ID: listID, UserID: userID}).Where(models.ListItem{ID: id}).Find(&foundItem).Error; err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting list item", InternalError: err}
	}

	return &foundItem, nil
}

// AddUserListItem adds a list item
func (s *DefaultListsService) AddUserListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error) {
	foundList, err := s.listsRepo.FindByID(listID, userID)
	if err != nil {
		return 0, err
	}

	if foundList == nil {
		return 0, &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	i := models.ListItem{}
	i.ListID = listID
	i.FromDto(dto)
	if err := s.db.Create(&i).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list item", InternalError: err}
	}

	return i.ID, nil
}

// RemoveUserListItem removes a list item
func (s *DefaultListsService) RemoveUserListItem(id int32, listID int32, userID int32) error {
	foundList, err := s.listsRepo.FindByID(listID, userID)
	if err != nil {
		return err
	}

	if foundList == nil {
		return &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	if err := s.db.Where(models.ListItem{ID: id, ListID: listID}).Delete(models.ListItem{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user list item", InternalError: err}
	}
	return nil
}

// UpdateUserListItem updates a list item
func (s *DefaultListsService) UpdateUserListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error {
	foundItem, err := s.getListItem(id, listID, userID)
	if err != nil {
		return err
	}

	foundItem.FromDto(dto)

	if err := s.db.Save(&foundItem).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating list item", InternalError: err}
	}

	return nil
}
