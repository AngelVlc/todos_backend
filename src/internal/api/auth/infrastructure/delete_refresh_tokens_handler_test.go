package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func TestDeleteRefreshTokensHandlerValidations(t *testing.T) {
	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not an array of ids", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

		result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})
}

func TestDeleteRefreshTokensHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	ids := []int32{int32(1), int32(2)}
	body, _ := json.Marshal(ids)

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("DeleteRefreshTokensByID", ids).Return(fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error deleting the refresh tokens")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an ok result if the refresh tokens are delete", func(t *testing.T) {
		mockedRepo.On("DeleteRefreshTokensByID", ids).Return(nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

		results.CheckOkResult(t, result, http.StatusNoContent)
		mockedRepo.AssertExpectations(t)
	})
}
