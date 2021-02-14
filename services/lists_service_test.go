package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

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

	svc := NewDefaultListsService(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`) VALUES (?,?)")).
			WithArgs(l.Name, userID)
	}

	t.Run("AddUserList() should return an error if insert fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserList(userID, &listDto)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserList() should insert the new list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserList(userID, &listDto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("RemoveUserList() should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserList(listID, userID)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserList() should delete the user list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserList(listID, userID)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ?, `userId` = ? WHERE `lists`.`id` = ?")).
			WithArgs(l.Name, userID, listID)
	}

	t.Run("UpdateUserList() should return an error if the update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.UpdateUserList(listID, userID, &listDto)

		appErrors.CheckUnexpectedError(t, err, "Error updating list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserList() should update the list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`  WHERE `lists`.`id` = ? ORDER BY `lists`.`id` ASC LIMIT 1")).
			WithArgs(listID).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, l.Name, userID))

		err := svc.UpdateUserList(listID, userID, &listDto)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedGetListQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("GetSingleUserList() should return an error if the query fails", func(t *testing.T) {
		dto := dtos.GetSingleListResultDto{}

		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetSingleUserList(listID, userID, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetSingleUserList() should get a single list", func(t *testing.T) {
		dto := dtos.GetSingleListResultDto{}

		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(listID).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		err := svc.GetSingleUserList(listID, userID, &dto)

		assert.Equal(t, "list", dto.Name)
		assert.Equal(t, len(dto.ListItems), 1)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

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
		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		_, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	expectedListItemOpPreviousQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(listID)
	}
	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`title`,`description`) VALUES (?,?,?)")).
			WithArgs(i.ListID, i.Title, i.Description)
	}

	t.Run("AddUserListItem() should return an error if insert fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserListItem() should insert the new list item", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserListItem(i.ListID, userID, &listItemDto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE (`listItems`.`id` = ?) AND (`listItems`.`listId` = ?)")).
			WithArgs(i.ID, i.ListID)
	}

	t.Run("RemoveUserListItem() should return an error if delete fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should delete the user list item", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list", userID))
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserListItem(i.ID, i.ListID, userID)

		assert.Nil(t, err)

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
