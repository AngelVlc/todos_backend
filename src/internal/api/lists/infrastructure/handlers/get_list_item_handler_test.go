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
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Item_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
			"id":     "111",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemRecord{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := GetListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestGetListItemHandler_Returns_The_List(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
			"id":     "111",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	list := domain.ListItemRecord{ID: 111, ListID: 11, Title: "title", Description: "desc"}
	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemRecord{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()

	result := GetListItemHandler(httptest.NewRecorder(), request(), h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	listRes, isOk := okRes.Content.(*infrastructure.ListItemResponse)
	require.Equal(t, true, isOk, "should be a list item response")

	assert.Equal(t, int32(111), listRes.ID)
	assert.Equal(t, int32(11), listRes.ListID)
	assert.Equal(t, "title", listRes.Title)
	assert.Equal(t, "desc", listRes.Description)
	mockedRepo.AssertExpectations(t)
}
