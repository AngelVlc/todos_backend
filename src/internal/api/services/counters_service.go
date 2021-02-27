package services

import (
	"fmt"

	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type CountersService interface {
	CreateCounterIfNotExists(name string) error
	IncrementCounter(name string) (int32, error)
}

type MockedCountersService struct {
	mock.Mock
}

func NewMockedCountersService() *MockedCountersService {
	return &MockedCountersService{}
}

func (m *MockedCountersService) CreateCounterIfNotExists(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockedCountersService) IncrementCounter(name string) (int32, error) {
	args := m.Called(name)
	return args.Get(0).(int32), args.Error(1)
}

type DefaultCountersService struct {
	db *gorm.DB
}

func NewDefaultCountersService(db *gorm.DB) *DefaultCountersService {
	return &DefaultCountersService{db}
}

func (s *DefaultCountersService) CreateCounterIfNotExists(name string) error {
	var counter models.Counter
	return s.db.Where(models.Counter{Name: name}).Attrs(models.Counter{Value: 0}).FirstOrCreate(&counter).Error
}

func (s *DefaultCountersService) IncrementCounter(name string) (int32, error) {
	var counter models.Counter
	err := s.db.Where(models.Counter{Name: name}).First(&counter).Error
	if err != nil {
		return -1, fmt.Errorf("error getting '%v' counter: %v", name, err)
	}
	counter.Value++
	err = s.db.Save(&counter).Error
	if err != nil {
		return -1, fmt.Errorf("error saving new '%v' counter value: %v", name, err)
	}
	return counter.Value, nil
}
