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
	listColumns := []string{"id", "name", "userId"}
	listItemsColumns := []string{"id", "listId", "title", "description"}

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	svc := NewDefaultListsService(db)

	t.Run("AddUserList() should return an error if insert fails", func(t *testing.T) {
		u := int32(11)
		l := models.List{
			Name: "list",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`) VALUES (?,?)")).
			WithArgs(l.Name, u).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserList(u, &l)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddUserList() should insert the new list", func(t *testing.T) {
		u := int32(11)
		l := models.List{
			Name: "list",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`) VALUES (?,?)")).
			WithArgs(l.Name, u).
			WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserList(u, &l)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUserList() should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, 22).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserList(11, 22)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUserList() should delete the user list", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, 22).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserList(11, 22)

		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUserList() should return an error if delete fails", func(t *testing.T) {
		u := int32(11)
		l := models.List{
			Name: "list",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ?, `userId` = ? WHERE `lists`.`id` = ?")).
			WithArgs(l.Name, u, 11).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.UpdateUserList(11, u, &l)

		appErrors.CheckUnexpectedError(t, err, "Error updating list", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUserList() should update the list", func(t *testing.T) {
		u := int32(11)
		l := models.List{
			Name: "list",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ?, `userId` = ? WHERE `lists`.`id` = ?")).
			WithArgs(l.Name, u, 11).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`  WHERE `lists`.`id` = ? ORDER BY `lists`.`id` ASC LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, l.Name, u))

		err := svc.UpdateUserList(11, u, &l)

		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetSingleUserList() should return an error if the query fails", func(t *testing.T) {
		u := int32(11)
		dto := dtos.GetSingleListResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, u).
			WillReturnError(fmt.Errorf("some error"))

		err := svc.GetSingleUserList(11, u, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetSingleUserList() should get a single list", func(t *testing.T) {
		u := int32(11)
		dto := dtos.GetSingleListResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, u).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", u))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(22, 11, "title", "description"))

		err := svc.GetSingleUserList(11, u, &dto)

		assert.Equal(t, "list", dto.Name)
		assert.Equal(t, len(dto.ListItems), 1)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetUserLists() should return an error if the query fails", func(t *testing.T) {
		u := int32(11)
		dto := []dtos.GetListsResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(u).
			WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUserLists(u, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user lists", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetUserLists() should return the user lists", func(t *testing.T) {
		u := int32(11)
		dto := []dtos.GetListsResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(u).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list1", u).AddRow(12, "list2", u))

		err := svc.GetUserLists(u, &dto)

		assert.Equal(t, len(dto), 2)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
