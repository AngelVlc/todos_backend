//+build !e2e

package reqcountermdw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sharedDomain "github.com/AngelVlc/todos/internal/api/shared/domain"
	sharedInfra "github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestRequestCounterMiddleware(t *testing.T) {
	mockedCountersRepo := sharedInfra.MockedCountersRepository{}
	md := RequestCounterMiddleware{&mockedCountersRepo}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	t.Run("should return 500 if incrementing the counter failt", func(t *testing.T) {
		mockedCountersRepo.On("FindByName", "requests").Return(nil, fmt.Errorf("some error")).Once()

		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "Error incrementing requests counter\n", string(response.Body.String()))
		mockedCountersRepo.AssertExpectations(t)
	})

	t.Run("should call next handler when the counter is incremented", func(t *testing.T) {
		counter := sharedDomain.Counter{}
		mockedCountersRepo.On("FindByName", "requests").Return(&counter, nil).Once()
		mockedCountersRepo.On("Update", &counter).Return(nil).Once()

		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
		mockedCountersRepo.AssertExpectations(t)
	})
}
