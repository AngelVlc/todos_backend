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

func TestListItemsServiceGetItem(t *testing.T) {
	mockedListItemsRepo := repositories.NewMockedListItemsRepository()

	svc := NewDefaultListItemsService(mockedListItemsRepo, nil)

	itemID := int32(111)
	listID := int32(11)
	userID := int32(1)

	t.Run("should return an error if the query fails", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.GetListItem(itemID, listID, userID)

		assert.Nil(t, dto)
		assert.Error(t, err)

		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the item doesn't exist", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(nil, nil).Once()

		dto, err := svc.GetListItem(itemID, listID, userID)

		assert.Nil(t, dto)
		assert.Nil(t, err)

		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should get an item", func(t *testing.T) {
		foundItem := models.ListItem{
			ID:          itemID,
			ListID:      listID,
			Title:       "title",
			Description: "desc",
		}

		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(&foundItem, nil).Once()

		dto, err := svc.GetListItem(itemID, listID, userID)

		assert.Nil(t, err)
		require.NotNil(t, dto)
		assert.Equal(t, "title", dto.Title)
		assert.Equal(t, "desc", dto.Description)

		mockedListItemsRepo.AssertExpectations(t)
	})
}

func TestListItemsServiceAddItem(t *testing.T) {
	mockedListItemsRepo := repositories.NewMockedListItemsRepository()
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListItemsService(mockedListItemsRepo, mockedListsRepo)

	listID := int32(11)
	userID := int32(1)

	listItemDto := dtos.ListItemDto{Title: "title", Description: "desc"}

	t.Run("should return an error if getting the list fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddListItem(listID, userID, &listItemDto)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the list doesn't exist", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, nil).Once()

		_, err := svc.AddListItem(listID, userID, &listItemDto)

		appErrors.CheckBadRequestError(t, err, "The list does not exist", "")

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if insert fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()
		mockedListItemsRepo.On("Insert", &models.ListItem{ListID: listID, Title: "title", Description: "desc"}).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddListItem(listID, userID, &listItemDto)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should insert the new list item", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()
		mockedListItemsRepo.On("Insert", &models.ListItem{ListID: listID, Title: "title", Description: "desc"}).Return(int32(111), nil).Once()

		id, err := svc.AddListItem(listID, userID, &listItemDto)

		assert.Equal(t, int32(111), id)
		assert.Nil(t, err)

		mockedListsRepo.AssertExpectations(t)
		mockedListItemsRepo.AssertExpectations(t)
	})
}

func TestListItemsServiceRemoveListItem(t *testing.T) {
	mockedListItemsRepo := repositories.NewMockedListItemsRepository()
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListItemsService(mockedListItemsRepo, mockedListsRepo)

	itemID := int32(111)
	listID := int32(11)
	userID := int32(1)

	t.Run("should return an error if getting the list fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.RemoveListItem(itemID, listID, userID)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the list doesn't exist", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(nil, nil).Once()

		err := svc.RemoveListItem(itemID, listID, userID)

		appErrors.CheckBadRequestError(t, err, "The list does not exist", "")

		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()
		mockedListItemsRepo.On("Remove", itemID, listID, userID).Return(fmt.Errorf("some error")).Once()

		err := svc.RemoveListItem(itemID, listID, userID)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should delete the user list item", func(t *testing.T) {
		mockedListsRepo.On("FindByID", listID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()
		mockedListItemsRepo.On("Remove", itemID, listID, userID).Return(nil).Once()

		err := svc.RemoveListItem(itemID, listID, userID)

		assert.Nil(t, err)

		mockedListsRepo.AssertExpectations(t)
		mockedListItemsRepo.AssertExpectations(t)
	})
}

func TestListItemsServiceUpdateListItem(t *testing.T) {
	mockedListItemsRepo := repositories.NewMockedListItemsRepository()

	svc := NewDefaultListItemsService(mockedListItemsRepo, nil)

	itemID := int32(111)
	listID := int32(11)
	userID := int32(1)
	listItemDto := dtos.ListItemDto{Title: "title", Description: "desc"}

	t.Run("should return an error if getting the list fails", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.UpdateListItem(itemID, listID, userID, &listItemDto)

		assert.Error(t, err)
		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the list doesn't exist", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(nil, nil).Once()

		err := svc.UpdateListItem(itemID, listID, userID, &listItemDto)

		appErrors.CheckBadRequestError(t, err, "The item does not exist", "")
		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(&models.ListItem{ID: itemID, Title: "ori title", Description: "ori desc", ListID: listID}, nil).Once()
		mockedListItemsRepo.On("Update", &models.ListItem{ID: itemID, Title: "title", Description: "desc", ListID: listID}).Return(fmt.Errorf("some error")).Once()

		err := svc.UpdateListItem(itemID, listID, userID, &listItemDto)

		assert.Error(t, err)
		mockedListItemsRepo.AssertExpectations(t)
	})

	t.Run("should update the item", func(t *testing.T) {
		mockedListItemsRepo.On("FindByID", itemID, listID, userID).Return(&models.ListItem{ID: itemID, Title: "ori title", Description: "ori desc", ListID: listID}, nil).Once()
		mockedListItemsRepo.On("Update", &models.ListItem{ID: itemID, Title: "title", Description: "desc", ListID: listID}).Return(nil).Once()

		err := svc.UpdateListItem(itemID, listID, userID, &listItemDto)

		assert.Nil(t, err)
		mockedListItemsRepo.AssertExpectations(t)
	})
}
