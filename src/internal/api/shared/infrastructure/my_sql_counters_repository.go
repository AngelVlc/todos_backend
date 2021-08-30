package infrastructure

import (
	"errors"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"gorm.io/gorm"
)

type MySqlCountersRepository struct {
	db *gorm.DB
}

func NewMySqlCountersRepository(db *gorm.DB) *MySqlCountersRepository {
	return &MySqlCountersRepository{db}
}

func (r *MySqlCountersRepository) FindByName(name string) (*domain.Counter, error) {
	foundCounter := domain.Counter{}
	err := r.db.Where(domain.Counter{Name: name}).First(&foundCounter).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &foundCounter, nil
}

func (r *MySqlCountersRepository) Create(counter *domain.Counter) error {
	return r.db.Create(counter).Error
}

func (r *MySqlCountersRepository) Update(counter *domain.Counter) error {
	return r.db.Save(&counter).Error
}
