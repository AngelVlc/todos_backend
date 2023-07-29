//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllRefreshTokensHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_Fails(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	paginatedOptions := sharedDomain.NewPaginationInfo(10, 0, "id", sharedDomain.OrderAsc)
	mockedRepo.On("GetAllRefreshTokens", request.Context(), paginatedOptions).Return(nil, fmt.Errorf("some error")).Once()

	result := GetAllRefreshTokensHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting refresh tokens")
	mockedRepo.AssertExpectations(t)
}

func TestGetAllRefreshTokensHandler_Returns_The_RefreshTokens(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	time1 := time.Now()
	time2 := time1.Add(1 * time.Hour)
	found := []domain.RefreshTokenRecord{
		{ID: 2, UserID: 1, ExpirationDate: time1},
		{ID: 5, UserID: 3, ExpirationDate: time2},
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	paginatedOptions := sharedDomain.NewPaginationInfo(10, 0, "id", sharedDomain.OrderAsc)
	mockedRepo.On("GetAllRefreshTokens", request.Context(), paginatedOptions).Return(found, nil)
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
}
