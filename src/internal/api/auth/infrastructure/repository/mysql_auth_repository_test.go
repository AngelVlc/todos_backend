//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySqlAuthRepository_ExistsRefreshToken_Returns_An_Error_If_The_Query_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	rt := "rt"
	userID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta("ELECT count(*) FROM `refresh_tokens` WHERE `refresh_tokens`.`userId` = ? AND `refresh_tokens`.`refreshToken` = ?")).
		WithArgs(userID, rt).
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.ExistsRefreshToken(context.Background(), domain.RefreshTokenEntity{RefreshToken: rt, UserID: userID})

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_FindRefreshTokenForUser_Returns_A_RefreshToken_If_Exists(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	rt := "rt"
	userID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta("ELECT count(*) FROM `refresh_tokens` WHERE `refresh_tokens`.`userId` = ? AND `refresh_tokens`.`refreshToken` = ?")).
		WithArgs(userID, rt).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	res, err := repo.ExistsRefreshToken(context.Background(), domain.RefreshTokenEntity{RefreshToken: rt, UserID: userID})

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_CreateRefreshTokenIfNotExist_Returns_An_Error_If_Creating_The_New_RefreshToken_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expDate, _ := time.Parse("2021-Jan-01", "2014-Feb-04")
	rt := domain.RefreshTokenEntity{UserID: 1, RefreshToken: "rt", ExpirationDate: expDate}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `refresh_tokens` (`userId`,`refreshToken`,`expirationDate`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
		WithArgs(rt.UserID, rt.RefreshToken, rt.ExpirationDate).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.CreateRefreshTokenIfNotExist(context.Background(), &rt)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_CreateRefreshTokenIfNotExist_Creates_A_New_RefreshToken(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expDate, _ := time.Parse("2021-Jan-01", "2014-Feb-04")
	rt := domain.RefreshTokenEntity{UserID: 1, RefreshToken: "rt", ExpirationDate: expDate}

	result := sqlmock.NewResult(12, 1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `refresh_tokens` (`userId`,`refreshToken`,`expirationDate`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
		WithArgs(rt.UserID, rt.RefreshToken, rt.ExpirationDate).
		WillReturnResult(result)
	mock.ExpectCommit()

	err := repo.CreateRefreshTokenIfNotExist(context.Background(), &rt)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuth_DeleteExpiredRefreshTokens_Returns_An_Error_If_The_Delete_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE expirationDate <= ?")).
		WithArgs(now).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteExpiredRefreshTokens(context.Background(), now)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuth_DeleteExpiredRefreshTokens_Deletes_The_Expired_RefreshTokens(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE expirationDate <= ?")).
		WithArgs(now).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteExpiredRefreshTokens(context.Background(), now)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_GetAllRefreshTokens_Returns_An_Error_If_The_Query_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	paginationInfo := sharedDomain.NewPaginationInfo(10, 10, "expirationDate", sharedDomain.OrderAsc)

	mock.ExpectQuery("SELECT id,userId,expirationDate FROM `refresh_tokens` ORDER BY expirationDate asc LIMIT 10 OFFSET 10").
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllRefreshTokens(context.Background(), paginationInfo)

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_GetAllRefreshTokens_Returns_The_RefreshTokens(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	paginationInfo := sharedDomain.NewPaginationInfo(10, 10, "expirationDate", sharedDomain.OrderAsc)

	columns := []string{"id", "userId", "expirationDate"}
	now := time.Now()
	mock.ExpectQuery("SELECT id,userId,expirationDate FROM `refresh_tokens` ORDER BY expirationDate asc LIMIT 10 OFFSET 10").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(11, 1, now))

	res, err := repo.GetAllRefreshTokens(context.Background(), paginationInfo)

	assert.Nil(t, err)
	require.IsType(t, []*domain.RefreshTokenEntity{}, res)
	require.Equal(t, 1, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, int32(1), res[0].UserID)
	assert.Equal(t, now, res[0].ExpirationDate)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_DeleteRefreshTokensByID_Returns_An_Error_If_The_Delete_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	ids := []int32{1, 2, 3}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE `refresh_tokens`.`id` IN (?,?,?)")).
		WithArgs(1, 2, 3).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteRefreshTokensByID(context.Background(), ids)

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_DeleteRefreshTokensByID_Deletes_The_RefreshTokens(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	ids := []int32{1, 2, 3}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE `refresh_tokens`.`id` IN (?,?,?)")).
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteRefreshTokensByID(context.Background(), ids)

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}
