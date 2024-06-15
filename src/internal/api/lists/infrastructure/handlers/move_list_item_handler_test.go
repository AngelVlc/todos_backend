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
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func moveRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestMoveListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_Origin_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
	}

	request := moveRequest()

	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := MoveListItemHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestMoveListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_Destination_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
	}

	request := moveRequest()

	originList := domain.ListRecord{ID: 11, UserID: 1, Name: "origin list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(originList, nil).Once()
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 20, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := MoveListItemHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "The destination list does not exist")
	mockedRepo.AssertExpectations(t)
}

func TestMoveListItemHandler_Returns_An_Error_If_The_ListItem_Does_Not_Exist_In_The_Origin_List(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
	}

	request := moveRequest()

	originList := domain.ListRecord{ID: 11, UserID: 1, Name: "origin list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(originList, nil).Once()
	destinationList := domain.ListRecord{ID: 20, UserID: 1, Name: "destination list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 20, UserID: 1}).Return(destinationList, nil).Once()

	result := MoveListItemHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "An item with id 5 doesn't exist in the original list")
	mockedRepo.AssertExpectations(t)
}

func TestMoveListItemHandler_Returns_An_Error_If_Updating_The_Origin_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
	}

	request := moveRequest()

	originListItem := &domain.ListItemRecord{ID: 5}
	originList := domain.ListRecord{ID: 11, UserID: 1, Name: "origin list", Items: []*domain.ListItemRecord{originListItem}}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(originList, nil).Once()
	destinationList := domain.ListRecord{ID: 20, UserID: 1, Name: "destination list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 20, UserID: 1}).Return(destinationList, nil).Once()
	originList.Items = []*domain.ListItemRecord{}
	mockedRepo.On("UpdateList", request.Context(), &originList).Return(fmt.Errorf("some error")).Once()

	result := MoveListItemHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the original list")
	mockedRepo.AssertExpectations(t)
}

func TestMoveListItemHandler_Returns_An_Error_If_The_Update_Of_The_Destination_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
	}

	request := moveRequest()

	originListItem := &domain.ListItemRecord{ID: 5}
	originList := domain.ListRecord{ID: 11, UserID: 1, Name: "origin list", Items: []*domain.ListItemRecord{originListItem}}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(originList, nil).Once()
	destinationList := domain.ListRecord{ID: 20, UserID: 1, Name: "destination list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 20, UserID: 1}).Return(destinationList, nil).Once()
	originList.Items = []*domain.ListItemRecord{}
	mockedRepo.On("UpdateList", request.Context(), &originList).Return(nil).Once()
	destinationList.Items = []*domain.ListItemRecord{originListItem}
	mockedRepo.On("UpdateList", request.Context(), &destinationList).Return(fmt.Errorf("some error")).Once()

	result := MoveListItemHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the destination list")
	mockedRepo.AssertExpectations(t)
}

func TestMoveListItemHandler_Updates_The_Lists_An_And_Sends_Two_ListCreatedOrUpdated_Event(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.MoveListItemInput{OriginListItemID: 5, DestinationListID: 20},
		EventBus:        &mockedEventBus,
	}

	request := moveRequest()

	originListItem := &domain.ListItemRecord{ID: 5}
	originList := domain.ListRecord{ID: 11, UserID: 1, Name: "origin list", Items: []*domain.ListItemRecord{originListItem}}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(originList, nil).Once()
	destinationList := domain.ListRecord{ID: 20, UserID: 1, Name: "destination list"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 20, UserID: 1}).Return(destinationList, nil).Once()
	originList.Items = []*domain.ListItemRecord{}
	mockedRepo.On("UpdateList", request.Context(), &originList).Return(nil).Once()
	destinationList.Items = []*domain.ListItemRecord{originListItem}
	mockedRepo.On("UpdateList", request.Context(), &destinationList).Return(nil).Once()

	mockedEventBus.On("Publish", events.ListUpdated, int32(11))
	mockedEventBus.On("Publish", events.ListUpdated, int32(20))

	mockedEventBus.Wg.Add(2)
	result := MoveListItemHandler(httptest.NewRecorder(), request, h)
	mockedEventBus.Wg.Wait()

	results.CheckOkResult(t, result, http.StatusOK)
	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
