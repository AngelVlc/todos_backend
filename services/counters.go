package services

import (
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
)

type CountersService struct {
	db *gorm.DB
}

func NewCountersService(db *gorm.DB) CountersService {
	return CountersService{db}
}

func (s *CountersService) CreateCounterIfNotExists(name string) {
	var counter models.Counter
	s.db.Where(models.Counter{Name: name}).Attrs(models.Counter{Value: 0}).FirstOrCreate(&counter)
}
