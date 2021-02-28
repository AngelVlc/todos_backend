//+build !e2e

package handlers

import (
	"os"
	"testing"

	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}

func CheckOkResult(t *testing.T, result interface{}, expectedStatus int) *results.OkResult {
	assert.NotNil(t, result)
	okRes, isOkResult := result.(results.OkResult)
	require.Equal(t, true, isOkResult, "should be an ok result")
	assert.Equal(t, expectedStatus, okRes.StatusCode)
	return &okRes
}

func CheckBadRequestErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(results.ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	badReqErr, isBadReqErr := errorRes.Err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, errorMsg, badReqErr.Error())
}

func CheckUnexpectedErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(results.ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unexpErr, isUnexError := errorRes.Err.(*appErrors.UnexpectedError)
	require.Equal(t, true, isUnexError, "should be an unexpected error")
	assert.Equal(t, errorMsg, unexpErr.Error())
}

func CheckUnauthorizedErrorErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(results.ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unauthpErr, isUnauthError := errorRes.Err.(*appErrors.UnauthorizedError)
	require.Equal(t, true, isUnauthError, "should be an unauthorized error")
	assert.Equal(t, errorMsg, unauthpErr.Error())
}
