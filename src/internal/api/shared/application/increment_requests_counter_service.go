package application

import (
	"fmt"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
)

type IncrementRequestsCounterService struct {
	repo domain.CountersRepository
}

func NewIncrementRequestsCounterService(repo domain.CountersRepository) *IncrementRequestsCounterService {
	return &IncrementRequestsCounterService{repo}
}

func (s *IncrementRequestsCounterService) IncrementRequestsCounter() (int32, error) {
	foundCounter, err := s.repo.FindByName("requests")
	if err != nil {
		return -1, fmt.Errorf("error getting 'requests' counter: %v", err)
	}

	if foundCounter == nil {
		return -1, fmt.Errorf("'requests' counter does not exist")
	}

	foundCounter.Value++

	err = s.repo.Update(foundCounter)
	if err != nil {
		return -1, fmt.Errorf("error updating 'requests' counter: %v", err)
	}

	return foundCounter.Value, nil
}
