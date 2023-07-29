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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateListItemHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Find_The_ListItem_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id":     "111",
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemRecord{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListItemHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Update_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id":     "111",
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listItem := domain.ListItemRecord{ID: 1, ListID: 11, Title: domain.ItemTitleValueObject("title"), Description: domain.ItemDescriptionValueObject("desc"), UserID: int32(1)}
	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemRecord{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(&listItem, nil).Once()
	mockedRepo.On("UpdateListItem", request().Context(), &listItem).Return(fmt.Errorf("some error")).Once()

	result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the list item")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListItemHandler_Updates_The_ListItem(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id":     "111",
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listItem := domain.ListItemRecord{ID: 1, ListID: 11, Title: domain.ItemTitleValueObject("oriTitle"), Description: domain.ItemDescriptionValueObject("oriDesc"), UserID: int32(1)}
	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemRecord{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(&listItem, nil).Once()
	mockedRepo.On("UpdateListItem", request().Context(), &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListItemRecord)
		*arg = domain.ListItemRecord{Title: "title", Description: "desc"}
	})

	result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(infrastructure.ListItemResponse)
	require.Equal(t, true, isOk, "should be a ListItemResponse")
	assert.Equal(t, "title", res.Title)
	assert.Equal(t, "desc", res.Description)

	mockedRepo.AssertExpectations(t)
}
