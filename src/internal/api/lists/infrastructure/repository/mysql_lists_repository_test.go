//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	listColumns      = []string{"id", "name", "userId", "itemsCount"}
	listItemsColumns = []string{"id", "listId", "userId", "title", "description", "position"}
)

func TestMySqlListsRepository_FindList_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
		WithArgs(listID, userID).
		WillReturnError(fmt.Errorf("some error"))

	_, err := repo.FindList(context.Background(), domain.ListRecord{ID: listID, UserID: userID})

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_FindList_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows(listColumns).
			AddRow(listID, "list1", userID, 3))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`listId` = ? ORDER BY position ASC")).
		WithArgs(listID).
		WillReturnRows(sqlmock.NewRows(listItemsColumns).
			AddRow(21, listID, userID, "item1_title", "item1_desc", 0).
			AddRow(31, listID, userID, "item2_title", "item2_desc", 1))

	res, err := repo.FindList(context.Background(), domain.ListRecord{ID: listID, UserID: userID})

	require.NotNil(t, res)
	require.IsType(t, domain.ListRecord{}, res)
	assert.Equal(t, listID, res.ID)
	assert.Equal(t, "list1", res.Name)
	assert.Equal(t, userID, res.UserID)
	assert.Equal(t, int32(3), res.ItemsCount)
	assert.Equal(t, 2, len(res.Items))
	assert.Equal(t, int32(21), res.Items[0].ID)
	assert.Equal(t, listID, res.Items[0].ListID)
	assert.Equal(t, userID, res.Items[0].UserID)
	assert.Equal(t, "item1_title", res.Items[0].Title)
	assert.Equal(t, "item1_desc", res.Items[0].Description)
	assert.Equal(t, int32(0), res.Items[0].Position)
	assert.Equal(t, int32(31), res.Items[1].ID)
	assert.Equal(t, listID, res.Items[1].ListID)
	assert.Equal(t, userID, res.Items[1].UserID)
	assert.Equal(t, "item2_title", res.Items[1].Title)
	assert.Equal(t, "item2_desc", res.Items[1].Description)
	assert.Equal(t, int32(1), res.Items[1].Position)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_ExistsList_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `lists` WHERE `lists`.`name` = ? AND `lists`.`userId` = ?")).
		WithArgs("list name", userID).
		WillReturnError(fmt.Errorf("some error"))

	repo := NewMySqlListsRepository(db)
	list := domain.ListRecord{Name: "list name", UserID: userID}

	res, err := repo.ExistsList(context.Background(), list)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_ExistsList_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `lists` WHERE `lists`.`name` = ? AND `lists`.`userId` = ?")).
		WithArgs("list name", userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	repo := NewMySqlListsRepository(db)
	list := domain.ListRecord{Name: "list name", UserID: userID}

	res, err := repo.ExistsList(context.Background(), list)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllLists_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)

	repo := NewMySqlListsRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`")).
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllLists(context.Background())

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllLists_When_It_Does_Not_Fail_Including_Items(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)

	repo := NewMySqlListsRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`")).
		WillReturnRows(sqlmock.NewRows(listColumns).
			AddRow(11, "list1", 1, 3).
			AddRow(12, "list2", 2, 4))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`listId` IN (?,?) ORDER BY position ASC")).
		WithArgs(11, 12).
		WillReturnRows(sqlmock.NewRows(listItemsColumns).
			AddRow(21, 11, 1, "list1_item1_title", "list1_item1_desc", 0).
			AddRow(31, 11, 1, "list1_item2_title", "list1_item2_desc", 1).
			AddRow(41, 12, 2, "list2_item1_title", "list2_item1_desc", 0).
			AddRow(51, 12, 2, "list2_item2_title", "list2_item2_desc", 1))

	res, err := repo.GetAllLists(context.Background())

	assert.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, "list1", res[0].Name)
	assert.Equal(t, int32(1), res[0].UserID)
	assert.Equal(t, int32(3), res[0].ItemsCount)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, "list1_item1_title", res[0].Items[0].Title)
	assert.Equal(t, "list1_item1_desc", res[0].Items[0].Description)
	assert.Equal(t, "list2", res[1].Name)
	assert.Equal(t, int32(2), res[1].UserID)
	assert.Equal(t, int32(4), res[1].ItemsCount)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllListsForUser_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)

	repo := NewMySqlListsRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`userId` = ?")).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllListsForUser(context.Background(), userID)

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllListsForUser_When_It_Does_Not_Fail_Without_Including_Items(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)

	repo := NewMySqlListsRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`userId` = ?")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows(listColumns).
			AddRow(11, "list1", userID, 3).
			AddRow(12, "list2", userID, 4))

	res, err := repo.GetAllListsForUser(context.Background(), userID)

	assert.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, "list1", res[0].Name)
	assert.Equal(t, userID, res[0].UserID)
	assert.Equal(t, int32(3), res[0].ItemsCount)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, "list2", res[1].Name)
	assert.Equal(t, userID, res[1].UserID)
	assert.Equal(t, int32(4), res[1].ItemsCount)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateList_When_The_Create_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`,`categoryId`,`itemsCount`) VALUES (?,?,?,?)")).
		WithArgs("list1", 1, 2, 0).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	list := domain.ListRecord{UserID: 1, Name: "list1", CategoryID: &sql.NullInt32{Int32: 2, Valid: true}}
	repo := NewMySqlListsRepository(db)

	err := repo.CreateList(context.Background(), &list)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateList_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`,`categoryId`,`itemsCount`) VALUES (?,?,?,?)")).
		WithArgs("list1", 1, 2, 0).
		WillReturnResult(sqlmock.NewResult(12, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`,`position`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `listId`=VALUES(`listId`)")).
		WithArgs(0, 1, "item1 title", "item1 desc", 0).
		WillReturnResult(sqlmock.NewResult(12, 0))
	mock.ExpectCommit()

	list := domain.ListRecord{
		UserID:     1,
		Name:       "list1",
		CategoryID: &sql.NullInt32{Int32: 2, Valid: true},
		Items: []*domain.ListItemRecord{
			{UserID: 1, Title: "item1 title", Description: "item1 desc", Position: 0},
		},
	}
	repo := NewMySqlListsRepository(db)

	err := repo.CreateList(context.Background(), &list)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_Deleting_The_ListItems_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ?")).
		WithArgs(listID).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteList(context.Background(), domain.ListRecord{ID: listID})

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_Deleting_The_List_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ?")).
		WithArgs(listID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE `lists`.`id` = ?")).
		WithArgs(listID).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteList(context.Background(), domain.ListRecord{ID: listID})

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ?")).
		WithArgs(listID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE `lists`.`id` = ?")).
		WithArgs(listID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteList(context.Background(), domain.ListRecord{ID: listID})

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateList_When_The_Update_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name`=?,`userId`=? WHERE `id` = ?")).
		WithArgs("list1", 1, 11).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	list := domain.ListRecord{ID: 11, UserID: 1, Name: "list1"}

	err := repo.UpdateList(context.Background(), &list)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateList_To_Remove_The_Category(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name`=?,`userId`=?,`categoryId`=? WHERE `id` = ?")).
		WithArgs("list1", 1, nil, 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE listId = ?")).
		WithArgs(11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewMySqlListsRepository(db)
	list := domain.ListRecord{ID: 11, UserID: 1, Name: "list1", CategoryID: &sql.NullInt32{Valid: false}}

	err := repo.UpdateList(context.Background(), &list)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateList_When_The_Update_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name`=?,`userId`=?,`categoryId`=? WHERE `id` = ?")).
		WithArgs("list1", 1, 2, 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE listId = ?")).
		WithArgs(11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewMySqlListsRepository(db)

	list := domain.ListRecord{ID: 11, UserID: 1, Name: "list1", CategoryID: &sql.NullInt32{Int32: 2, Valid: true}}

	err := repo.UpdateList(context.Background(), &list)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_IncrementListCounter_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=(SELECT COUNT(id) FROM `listItems` WHERE `listItems`.`listId` = ?) WHERE `lists`.`id` = ?")).
		WithArgs(11, 11).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.UpdateListItemsCount(context.Background(), 11)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateListItemsCounter_When_The_Update_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=(SELECT COUNT(id) FROM `listItems` WHERE `listItems`.`listId` = ?) WHERE `lists`.`id` = ?")).
		WithArgs(11, 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateListItemsCount(context.Background(), 11)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}
