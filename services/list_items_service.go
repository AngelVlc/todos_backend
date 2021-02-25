package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/stretchr/testify/mock"
)

type ListItemsService interface {
	GetListItem(id int32, listID int32, userID int32) (*dtos.ListItemResponseDto, error)
	AddListItem(listID int32, userID int32, dto *dtos.ListItemDto) (int32, error)
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

	return s.itemsRepo.Insert(&i)
}
