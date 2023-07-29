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

func TestUpdateListHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Is_Trying_To_Update_The_List_Name_But_The_Query_To_Check_The_Name_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	oldListName, _ := domain.NewListNameValueObject("oldName")
	list := domain.ListRecord{ID: int32(11), Name: oldListName, UserID: int32(11)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("ExistsList", request().Context(), &domain.ListRecord{Name: listName, UserID: int32(1)}).Return(false, fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Updating_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	list := domain.ListRecord{ID: int32(11), Name: listName, UserID: int32(11)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("UpdateList", request().Context(), &list).Return(fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the user list")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Getting_The_List_Items_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	originalListName, _ := domain.NewListNameValueObject("originalName")
	newListName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{ID: int32(11), Name: originalListName, UserID: int32(1)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("ExistsList", request().Context(), &domain.ListRecord{Name: newListName, UserID: int32(1)}).Return(false, nil).Once()
	mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListRecord)
		*arg = domain.ListRecord{Name: newListName}
	})
	mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting all list items")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Bulk_Update_Of_Their_Items_Position_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	originalListName, _ := domain.NewListNameValueObject("originalName")
	newListName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{ID: int32(11), Name: originalListName, UserID: int32(1)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("ExistsList", request().Context(), &domain.ListRecord{Name: newListName, UserID: int32(1)}).Return(false, nil).Once()
	mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListRecord)
		*arg = domain.ListRecord{Name: newListName}
	})
	listItems := []domain.ListItemRecord{{ID: int32(1), Position: int32(0)}, {ID: int32(2), Position: int32(1)}}
	mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(listItems, nil).Once()
	mockedRepo.On("BulkUpdateListItems", request().Context(), listItems).Return(fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error bulk updating")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Updates_The_List_And_The_Position_Of_Their_Items(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &domain.UpdateListInput{Name: listName, IDsByPosition: []int32{int32(2), int32(1)}},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	originalListName, _ := domain.NewListNameValueObject("originalName")
	newListName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{ID: int32(11), Name: originalListName, UserID: int32(1)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("ExistsList", request().Context(), &domain.ListRecord{Name: newListName, UserID: int32(1)}).Return(false, nil).Once()
	mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListRecord)
		*arg = domain.ListRecord{Name: newListName}
	})
	listItems := []domain.ListItemRecord{{ID: int32(1), Position: int32(0)}, {ID: int32(2), Position: int32(1)}}
	mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(listItems, nil).Once()
	listItems = []domain.ListItemRecord{{ID: int32(1), Position: int32(1)}, {ID: int32(2), Position: int32(0)}}
	mockedRepo.On("BulkUpdateListItems", request().Context(), listItems).Return(nil).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request(), h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(infrastructure.ListResponse)
	require.True(t, isOk, "should be a ListResponse")
	assert.Equal(t, "list1", res.Name)

	mockedRepo.AssertExpectations(t)
}
