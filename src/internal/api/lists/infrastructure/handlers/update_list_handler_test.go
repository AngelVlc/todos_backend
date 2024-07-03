//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestUpdateListHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	request := updateRequest()

	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Is_Trying_To_Update_The_List_Name_But_The_Query_To_Check_If_The_A_List_With_The_Same_Name_Exists_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	request := updateRequest()

	foundList := domain.ListRecord{ID: 11, Name: "oldName", UserID: 11}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&foundList, nil).Once()
	mockedRepo.On("ExistsList", request.Context(), domain.ListRecord{Name: "list1", UserID: 1}).Return(false, fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "Error checking if a list with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Updating_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	request := updateRequest()

	foundList := domain.ListRecord{
		ID:         int32(11),
		Name:       "list1",
		UserID:     1,
		Items:      []domain.ListItemRecord{},
		CategoryID: &sql.NullInt32{Valid: false},
	}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&foundList, nil).Once()
	mockedRepo.On("UpdateList", request.Context(), &foundList).Return(fmt.Errorf("some error")).Once()

	result := UpdateListHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the user list")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateListHandler_Updates_The_List_Name_And_Sends_The_ListCreatedOrUpdated_Event(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	newListName, _ := domain.NewListNameValueObject("list new name")
	newCategoryID := int32(5)
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: newListName, CategoryID: &newCategoryID},
		EventBus:        &mockedEventBus,
	}

	request := updateRequest()

	recordToUpdate := domain.ListRecord{
		ID:         int32(11),
		Name:       "list new name",
		UserID:     1,
		Items:      []domain.ListItemRecord{},
		CategoryID: &sql.NullInt32{Int32: 5, Valid: true},
	}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&domain.ListRecord{Name: "list1"}, nil).Once()
	mockedRepo.On("ExistsList", request.Context(), domain.ListRecord{Name: "list new name", UserID: 1}).Return(false, nil).Once()
	mockedRepo.On("UpdateList", request.Context(), &recordToUpdate).Return(nil).Once()

	mockedEventBus.On("Publish", events.ListUpdated, int32(11))

	mockedEventBus.Wg.Add(1)
	result := UpdateListHandler(httptest.NewRecorder(), request, h)
	mockedEventBus.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(*domain.ListEntity)
	require.True(t, isOk, "should be a ListEntity")
	assert.Equal(t, "list new name", res.Name.String())
	assert.Equal(t, int32(5), *res.CategoryID)

	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
