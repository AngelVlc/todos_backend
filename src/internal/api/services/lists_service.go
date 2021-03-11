package services

import (
	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/AngelVlc/todos/internal/api/repositories"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsService interface {
	AddUserList(userID int32, dto *dtos.ListDto) (int32, error)
	RemoveUserList(id int32, userID int32) error
	UpdateUserList(id int32, userID int32, dto *dtos.ListDto) error
	GetUserList(id int32, userID int32) (*dtos.ListResponseDto, error)
	GetUserLists(userID int32) ([]*dtos.ListResponseDto, error)
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

	return s.listsRepo.Create(&l)
}

// RemoveUserList removes a list
func (s *DefaultListsService) RemoveUserList(id int32, userID int32) error {
	return s.listsRepo.Delete(id, userID)
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
