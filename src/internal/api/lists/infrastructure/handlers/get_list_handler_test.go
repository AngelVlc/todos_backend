//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestGetListHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	request := getRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := GetListHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestGetListHandler_Returns_The_List(t *testing.T) {
	request := getRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	foundList := domain.ListRecord{ID: 11, Name: "list1", ItemsCount: 4}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&foundList, nil).Once()

	result := GetListHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	listRes, isOk := okRes.Content.(*domain.ListEntity)
	require.Equal(t, true, isOk, "should be a ListEntity")

	assert.Equal(t, int32(11), listRes.ID)
	assert.Equal(t, "list1", listRes.Name.String())
	assert.Equal(t, int32(4), listRes.ItemsCount)
	mockedRepo.AssertExpectations(t)
}
