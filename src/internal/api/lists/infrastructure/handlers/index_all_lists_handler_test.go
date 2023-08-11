//go:build !e2e
// +build !e2e

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func TestIndexAllListsHandler(t *testing.T) {
	mockedEventBus := events.MockedEventBus{}

	h := handler.Handler{
		EventBus: &mockedEventBus,
	}

	mockedEventBus.On("Publish", events.IndexAllListsRequested, nil)

	request, _ := http.NewRequest(http.MethodPost, "/", nil)

	mockedEventBus.Wg.Add(1)
	result := IndexAllListsHandler(httptest.NewRecorder(), request, h)
	mockedEventBus.Wg.Wait()

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedEventBus.AssertExpectations(t)
}
