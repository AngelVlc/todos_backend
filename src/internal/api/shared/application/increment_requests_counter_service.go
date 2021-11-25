package application

import (
	"fmt"
	"sync"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
)

type IncrementRequestsCounterService struct {
	repo domain.CountersRepository
	mu   sync.Mutex
}

func NewIncrementRequestsCounterService(repo domain.CountersRepository) *IncrementRequestsCounterService {
	return &IncrementRequestsCounterService{repo, sync.Mutex{}}
}

func (s *IncrementRequestsCounterService) IncrementRequestsCounter() (int32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
