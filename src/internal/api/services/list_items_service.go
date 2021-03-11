package services

import (
	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/AngelVlc/todos/internal/api/repositories"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/mock"
)

type ListItemsService interface {
	GetListItem(id int32, listID int32, userID int32) (*dtos.ListItemResponseDto, error)
	AddListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error)
	RemoveListItem(id int32, listID int32, userID int32) error
	UpdateListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error
}

type MockedListItemsService struct {
	mock.Mock
}

func NewMockedListItemsService() *MockedListItemsService {
	return &MockedListItemsService{}
}

func (m *MockedListItemsService) GetListItem(id int32, listID int32, userID int32) (*dtos.ListItemResponseDto, error) {
	args := m.Called(id, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.ListItemResponseDto), args.Error(1)
}

func (m *MockedListItemsService) AddListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error) {
	args := m.Called(listID, userID, dto)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListItemsService) RemoveListItem(id int32, listID int32, userID int32) error {
	args := m.Called(id, listID, userID)
	return args.Error(0)
}

func (m *MockedListItemsService) UpdateListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error {
	args := m.Called(id, listID, userID, dto)
	return args.Error(0)
}

type DefaultListItemsService struct {
	itemsRepo repositories.ListItemsRepository
	listsRepo repositories.ListsRepository
}

func NewDefaultListItemsService(itemsRepo repositories.ListItemsRepository, listsRepo repositories.ListsRepository) *DefaultListItemsService {
	return &DefaultListItemsService{itemsRepo, listsRepo}
}

func (s *DefaultListItemsService) GetListItem(id int32, listID int32, userID int32) (*dtos.ListItemResponseDto, error) {
	foundItem, err := s.itemsRepo.FindByID(id, listID, userID)
	if err != nil {
		return nil, err
	}

	if foundItem == nil {
		return nil, nil
	}

	return foundItem.ToResponseDto(), nil
}

func (s *DefaultListItemsService) AddListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error) {
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

	return s.itemsRepo.Create(&i)
}

// RemoveListItem removes an item
func (s *DefaultListItemsService) RemoveListItem(id int32, listID int32, userID int32) error {
	foundList, err := s.listsRepo.FindByID(listID, userID)
	if err != nil {
		return err
	}

	if foundList == nil {
		return &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	return s.itemsRepo.Delete(id, listID, userID)
}

// UpdateListItem updates an item
func (s *DefaultListItemsService) UpdateListItem(id int32, listID int32, userID int32, dto *dtos.ListItemDto) error {
	foundItem, err := s.itemsRepo.FindByID(id, listID, userID)
	if err != nil {
		return err
	}

	if foundItem == nil {
		return &appErrors.BadRequestError{Msg: "The item does not exist"}
	}

	foundItem.FromDto(dto)

	return s.itemsRepo.Update(foundItem)
}
