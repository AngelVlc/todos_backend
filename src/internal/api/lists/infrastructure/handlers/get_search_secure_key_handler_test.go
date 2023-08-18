//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getKeyRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestGetSearchSecureKeyHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Check_If_A_List_With_Generating_The_Key_Fails(t *testing.T) {
	request := getKeyRequest()

	mockedSearchClient := search.MockedSearchIndexClient{}

	h := handler.Handler{
		SearchClient: &mockedSearchClient,
	}

	mockedSearchClient.On("GenerateSecuredApiKey", "userID:1").Return("", fmt.Errorf("some error")).Once()

	result := GetSearchSecureKeyHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting the search key")
	mockedSearchClient.AssertExpectations(t)
}

func TestGetSearchSecureKeyHandler_Returns_The_Search_Key(t *testing.T) {
	request := getKeyRequest()

	mockedSearchClient := search.MockedSearchIndexClient{}

	h := handler.Handler{
		SearchClient: &mockedSearchClient,
	}

	mockedSearchClient.On("GenerateSecuredApiKey", "userID:1").Return("searchKey", nil).Once()

	result := GetSearchSecureKeyHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(string)
	require.True(t, isOk, "should be a string")
	assert.Equal(t, "searchKey", res)
	mockedSearchClient.AssertExpectations(t)
}
