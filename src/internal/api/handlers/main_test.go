package handlers

import (
	"os"
	"testing"

	appErrors "github.com/AngelVlc/todos/internal/api/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}

func CheckOkResult(t *testing.T, result interface{}, expectedStatus int) *okResult {
	assert.NotNil(t, result)
	okRes, isOkResult := result.(okResult)
	require.Equal(t, true, isOkResult, "should be an ok result")
	assert.Equal(t, expectedStatus, okRes.statusCode)
	return &okRes
}

func CheckBadRequestErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	badReqErr, isBadReqErr := errorRes.err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, errorMsg, badReqErr.Error())
}

func CheckUnexpectedErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unexpErr, isUnexError := errorRes.err.(*appErrors.UnexpectedError)
	require.Equal(t, true, isUnexError, "should be an unexpected error")
	assert.Equal(t, errorMsg, unexpErr.Error())
}

func CheckUnauthorizedErrorErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unauthpErr, isUnauthError := errorRes.err.(*appErrors.UnauthorizedError)
	require.Equal(t, true, isUnauthError, "should be an unauthorized error")
	assert.Equal(t, errorMsg, unauthpErr.Error())
}
