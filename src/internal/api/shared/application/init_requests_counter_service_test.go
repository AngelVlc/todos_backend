//+build !e2e

package application

import (
	"fmt"
	"testing"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestInitRequestsCounterService(t *testing.T) {
	mockedRepo := infrastructure.MockedCountersRepository{}

	svc := NewInitRequestsCounterService(&mockedRepo)

	t.Run("should return an error if finding the counter fails", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(nil, fmt.Errorf("some error")).Once()

		err := svc.InitRequestsCounter()

		assert.EqualError(t, err, "error getting 'requests' counter: some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return nil if the counter already exists", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(&domain.Counter{}, nil).Once()

		err := svc.InitRequestsCounter()

		assert.Nil(t, err)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an error if creating the counter fails", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(nil, nil).Once()
		counter := domain.Counter{Name: "requests", Value: 0}
		mockedRepo.On("Create", &counter).Return(fmt.Errorf("some error")).Once()

		err := svc.InitRequestsCounter()

		assert.EqualError(t, err, "error creating 'requests' counter: some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return nil if the counter is created", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(nil, nil).Once()
		counter := domain.Counter{Name: "requests", Value: 0}
		mockedRepo.On("Create", &counter).Return(nil).Once()

		err := svc.InitRequestsCounter()

		assert.Nil(t, err)
		mockedRepo.AssertExpectations(t)
	})
}
