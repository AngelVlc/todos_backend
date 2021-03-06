//+build !e2e

package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllRefreshTokensHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	t.Run("Should return an error result with an unexpected error if the query fails", func(t *testing.T) {
		mockedRepo.On("GetAllRefreshTokens").Return(nil, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := GetAllRefreshTokensHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting refresh tokens")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return the users if the query does not fail", func(t *testing.T) {
		time1 := time.Now()
		time2 := time1.Add(1 * time.Hour)
		found := []domain.RefreshToken{
			{ID: 2, UserID: 1, ExpirationDate: time1},
			{ID: 5, UserID: 3, ExpirationDate: time2},
		}

		mockedRepo.On("GetAllRefreshTokens").Return(found, nil)
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		result := GetAllRefreshTokensHandler(httptest.NewRecorder(), request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		rtRes, isOk := okRes.Content.([]RefreshTokenResponse)
		require.Equal(t, true, isOk, "should be an array of refresh token response")

		require.Equal(t, len(rtRes), 2)
		assert.Equal(t, int32(2), rtRes[0].ID)
		assert.Equal(t, int32(1), rtRes[0].UserID)
		assert.Equal(t, time1, rtRes[0].ExpirationDate)
		assert.Equal(t, int32(5), rtRes[1].ID)
		assert.Equal(t, int32(3), rtRes[1].UserID)
		assert.Equal(t, time2, rtRes[1].ExpirationDate)

		mockedRepo.AssertExpectations(t)
	})
}
