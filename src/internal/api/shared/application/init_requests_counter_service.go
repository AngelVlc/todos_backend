package application

import (
	"fmt"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
)

type InitRequestsCounterService struct {
	repo domain.CountersRepository
}

func NewInitRequestsCounterService(repo domain.CountersRepository) *InitRequestsCounterService {
	return &InitRequestsCounterService{repo}
}

func (s *InitRequestsCounterService) InitRequestsCounter() error {
	foundCounter, err := s.repo.FindByName("requests")
	if err != nil {
		return fmt.Errorf("error getting 'requests' counter: %v", err)
	}

	if foundCounter != nil {
		return nil
	}

	newCounter := domain.Counter{Name: "requests", Value: 0}

	err = s.repo.Create(&newCounter)
	if err != nil {
		return fmt.Errorf("error creating 'requests' counter: %v", err)
	}

	return nil
}
