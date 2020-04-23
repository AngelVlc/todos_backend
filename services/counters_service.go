package services

import (
	"fmt"

	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
)

type CountersService struct {
	db *gorm.DB
}

func NewCountersService(db *gorm.DB) CountersService {
	return CountersService{db}
}

func (s *CountersService) CreateCounterIfNotExists(name string) error {
	var counter models.Counter
	return s.db.Where(models.Counter{Name: name}).Attrs(models.Counter{Value: 0}).FirstOrCreate(&counter).Error
}

func (s *CountersService) IncrementCounter(name string) (int32, error) {
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
