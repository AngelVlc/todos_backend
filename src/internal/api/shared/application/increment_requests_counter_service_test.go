//+build !e2e

package application

import (
	"fmt"
	"testing"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestIncementRequestsCounterService(t *testing.T) {
	mockedRepo := infrastructure.MockedCountersRepository{}

	svc := NewIncrementRequestsCounterService(&mockedRepo)

	t.Run("should return an error if finding the counter fails", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.IncrementRequestsCounter()

		assert.EqualError(t, err, "error getting 'requests' counter: some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the counter doesn't exist", func(t *testing.T) {
		mockedRepo.On("FindByName", "requests").Return(nil, nil).Once()

		_, err := svc.IncrementRequestsCounter()

		assert.EqualError(t, err, "'requests' counter does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should an error if the update fails", func(t *testing.T) {
		counter := domain.Counter{Value: 100}
		mockedRepo.On("FindByName", "requests").Return(&counter, nil).Once()
		counter.Value++
		mockedRepo.On("Update", &counter).Return(fmt.Errorf("some error")).Once()

		_, err := svc.IncrementRequestsCounter()

		assert.EqualError(t, err, "error updating 'requests' counter: some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should the incremented value after the update", func(t *testing.T) {
		counter := domain.Counter{Value: 100}
		mockedRepo.On("FindByName", "requests").Return(&counter, nil).Once()
		mockedRepo.On("Update", &counter).Return(nil).Once()

		v, err := svc.IncrementRequestsCounter()

		assert.Equal(t, int32(101), v)
		assert.Nil(t, err)
		mockedRepo.AssertExpectations(t)
	})
}
