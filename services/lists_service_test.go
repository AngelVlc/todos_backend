package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListsServiceAddList(t *testing.T) {
	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(nil, mockedListsRepo)

	userID := int32(1)

	listDto := dtos.ListDto{Name: "list1"}

	t.Run("should return an error if insert fails", func(t *testing.T) {
		mockedListsRepo.On("Insert", &models.List{Name: "list1", UserID: userID}).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUserList(userID, &listDto)

		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should insert the new list", func(t *testing.T) {
		mockedListsRepo.On("Insert", &models.List{Name: "list1", UserID: userID}).Return(int32(12), nil).Once()

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
		mockedListsRepo.On("Remove", listID, userID).Return(fmt.Errorf("some error")).Once()

		err := svc.RemoveUserList(listID, userID)

		assert.Error(t, err)
		mockedListsRepo.AssertExpectations(t)
	})

	t.Run("should delete the user list", func(t *testing.T) {
		mockedListsRepo.On("Remove", listID, userID).Return(nil).Once()

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

	t.Run("should return an error if the list don't exist", func(t *testing.T) {
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

func TestListsService(t *testing.T) {
	listColumns := []string{"id", "name", "userID"}
	listItemsColumns := []string{"id", "listId", "title", "description"}

	listID := int32(15)
	itemID := int32(5)
	userID := int32(11)

	listDto := dtos.ListDto{
		Name: "list",
	}

	l := models.List{}
	l.FromDto(&listDto)

	i := models.ListItem{
		ID:          itemID,
		ListID:      listID,
		Title:       "title",
		Description: "description",
	}

	listItemDto := dtos.ListItemDto{
		Title:       i.Title,
		Description: i.Description,
	}

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	mockedListsRepo := repositories.NewMockedListsRepository()

	svc := NewDefaultListsService(db, mockedListsRepo)

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(userID)
	}

	t.Run("GetUserLists() should return an error if the query fails", func(t *testing.T) {
		dto := []dtos.GetListsResultDto{}

		expectedGetListsQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUserLists(userID, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user lists", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetUserLists() should return the user lists", func(t *testing.T) {
		dto := []dtos.GetListsResultDto{}

		expectedGetListsQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list1", userID).AddRow(12, "list2", userID))

		err := svc.GetUserLists(userID, &dto)

		assert.Equal(t, len(dto), 2)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT `listItems`.* FROM `listItems` JOIN lists on listItems.listId=lists.id WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?) AND (`listItems`.`id` = ?)")).
			WithArgs(listID, userID, itemID)
	}

	t.Run("GetUserListItem() should return an error if the query fails", func(t *testing.T) {

		dto := dtos.GetItemResultDto{}

		expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUserListItem(itemID, listID, userID, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetUserListItem() should get a single item", func(t *testing.T) {
		dto := dtos.GetItemResultDto{}

		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		err := svc.GetUserListItem(itemID, listID, userID, &dto)

		assert.Equal(t, "title", dto.Title)
		assert.Equal(t, "description", dto.Description)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`title`,`description`) VALUES (?,?,?)")).
			WithArgs(i.ListID, i.Title, i.Description)
	}

	t.Run("AddUserListItem() should return an error if insert fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()

		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list item", "some error")

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	t.Run("AddUserListItem() should insert the new list item", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()

		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		assert.Error(t, err)

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE (`listItems`.`id` = ?) AND (`listItems`.`listId` = ?)")).
			WithArgs(i.ID, i.ListID)
	}

	t.Run("RemoveUserListItem() should return an error if delete fails", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()

		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list item", "some error")

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should delete the user list item", func(t *testing.T) {
		mockedListsRepo.On("FindByID", i.ListID, userID).Return(&models.List{ID: listID, Name: "list1", UserID: userID}, nil).Once()

		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		assert.Nil(t, err)

		mockedListsRepo.AssertExpectations(t)
		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.UpdateUserListItem(i.ID, i.ListID, userID, &listItemDto)

		appErrors.CheckUnexpectedError(t, err, "Error getting list item", "some error")

		checkMockExpectations(t, mock)
	})

	expectedUpdateListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `listItems` SET `listId` = ?, `title` = ?, `description` = ? WHERE `listItems`.`id` = ?")).
			WithArgs(i.ListID, i.Title, i.Description, i.ID)
	}

	t.Run("UpdateUserListItem() should return an error if delete fails", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.UpdateUserListItem(i.ID, i.ListID, userID, &listItemDto)

		appErrors.CheckUnexpectedError(t, err, "Error updating list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserListItem() should update the list", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`id` = ? ORDER BY `listItems`.`id` ASC LIMIT 1")).
			WithArgs(itemID).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, l.Name, userID))

		err := svc.UpdateUserListItem(i.ID, i.ListID, userID, &listItemDto)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}
