package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func TestDeleteRefreshTokensHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Request_Does_Not_Have_Body(t *testing.T) {
	h := handler.Handler{}
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestDeleteRefreshTokensHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Body_Is_Not_An_Array_Of_Ids(t *testing.T) {
	h := handler.Handler{}
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

	result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestDeleteRefreshTokensHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Find_The_User_Fails(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	ids := []int32{int32(1), int32(2)}
	body, _ := json.Marshal(ids)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedRepo.On("DeleteRefreshTokensByID", request.Context(), ids).Return(fmt.Errorf("some error")).Once()

	result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the refresh tokens")
	mockedRepo.AssertExpectations(t)
}

func TestDeleteRefreshTokensHandler_Returns_An_Ok_Result_If_The_RefreshTokens_Are_Deleted(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	ids := []int32{int32(1), int32(2)}
	body, _ := json.Marshal(ids)

	request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockedRepo.On("DeleteRefreshTokensByID", request.Context(), ids).Return(nil).Once()

	result := DeleteRefreshTokensHandler(httptest.NewRecorder(), request, h)

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedRepo.AssertExpectations(t)
}
