package services

import (
	"fmt"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListsServiceAddList(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	userID := int32(1)

	listDto := dtos.ListDto{Name: "list1"}

	t.Run("should return an error if create fails", func(t *testing.T) {
		mockedListsRepo.On("Create", &models.List{Name: "list1", UserID: userID}).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUserList(userID, &listDto)

		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should create the new list", func(t *testing.T) {
		mockedListsRepo.On("Create", &models.List{Name: "list1", UserID: userID}).Return(int32(12), nil).Once()

		id, err := svc.AddUserList(userID, &listDto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)
		mockedListsRepo.AssertExpectations(t)
	})
}

func TestListsServiceRemoveList(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	userID := int32(1)
	listID := int32(11)

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mockedListsRepo.On("Delete", listID, userID).Return(fmt.Errorf("some error")).Once()

		err := svc.RemoveUserList(listID, userID)

		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should delete the user list", func(t *testing.T) {
		mockedListsRepo.On("Delete", listID, userID).Return(nil).Once()

		err := svc.RemoveUserList(listID, userID)

		assert.Nil(t, err)
		mockedListsRepo.AssertExpectations(t)
	})
}

func TestListsServiceUpdateList(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	listDto := dtos.ListDto{Name: "list1"}

	t.Run("should return an error if finding the list fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.UpdateUserList(11, 1, &listDto)

		require.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the list fails does not exist", func(t *testing.T) {
		mockedListsRepo.On("FindByID", int32(11), int32(1)).Return(nil, nil).Once()

		err := svc.UpdateUserList(11, 1, &listDto)

		appErrors.CheckBadRequestError(t, err, "The list does not exist", "")
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", int32(11), int32(1)).Return(&models.List{ID: int32(11), UserID: int32(1), Name: "ori"}, nil).Once()
		mockedListsRepo.On("Update", &models.List{ID: int32(11), UserID: int32(1), Name: "list1"}).Return(fmt.Errorf("some error")).Once()

		err := svc.UpdateUserList(11, 1, &listDto)

		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should update the list", func(t *testing.T) {
		mockedListsRepo.On("FindByID", int32(11), int32(1)).Return(&models.List{ID: int32(11), UserID: int32(1), Name: "ori"}, nil).Once()
		mockedListsRepo.On("Update", &models.List{ID: int32(11), UserID: int32(1), Name: "list1"}).Return(nil).Once()

		err := svc.UpdateUserList(11, 1, &listDto)

		assert.Nil(t, err)
		mockedListsRepo.AssertExpectations(t)
	})
}

func TestListsServiceGetList(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	listID := int32(11)
	userID := int32(1)

	t.Run("should return an error if the query fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.GetUserList(listID, userID)

		assert.Nil(t, dto)
		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the list doesn't exist", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, nil).Once()

		dto, err := svc.GetUserList(listID, userID)

		assert.Nil(t, dto)
		assert.Nil(t, err)

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should get a single list", func(t *testing.T) {
		foundList := models.List{
			ID:     listID,
			UserID: userID,
			Name:   "list1",
			ListItems: []*models.ListItem{
				{ID: int32(111), ListID: listID, Title: "title", Description: "desc"},
			},
		}

		mockedListsRepo.On("FindByID", listID, userID).Return(&foundList, nil).Once()

		dto, err := svc.GetUserList(listID, userID)

		assert.Nil(t, err)
		require.NotNil(t, dto)
		assert.Equal(t, "list1", dto.Name)
		require.Equal(t, 1, len(dto.ListItems))
		assert.Equal(t, "title", dto.ListItems[0].Title)
		assert.Equal(t, "desc", dto.ListItems[0].Description)

		mockedListsRepo.AssertExpectations(t)
	})
}

func TestListsServiceGetUserLists(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	userID := int32(1)

	t.Run("should return an error if the query fails", func(t *testing.T) {
		mockedListsRepo.On("GetAll", userID).Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.GetUserLists(userID)

		assert.Nil(t, dto)
		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return the user lists", func(t *testing.T) {
		found := []*models.List{
			{ID: 1, UserID: userID, Name: "list1"},
			{ID: 2, UserID: userID, Name: "list2"},
		}

		mockedListsRepo.On("GetAll", userID).Return(found, nil).Once()

		dto, err := svc.GetUserLists(userID)

		assert.Nil(t, err)
		require.NotNil(t, dto)
		require.Equal(t, 2, len(dto))
		assert.Equal(t, int32(1), dto[0].ID)
		assert.Equal(t, "list1", dto[0].Name)
		assert.Equal(t, int32(2), dto[1].ID)
		assert.Equal(t, "list2", dto[1].Name)
		mockedListsRepo.AssertExpectations(t)
	})
}
