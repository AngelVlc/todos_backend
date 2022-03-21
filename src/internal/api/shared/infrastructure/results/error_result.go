package results

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorResult struct {
	Err error
}

func (e ErrorResult) IsError() bool {
	return true
}

func CheckError(t *testing.T, result interface{}, errorMsg string) {
	require.NotNil(t, result)
	errorRes, isErrorResult := result.(ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")
	assert.Equal(t, errorMsg, errorRes.Err.Error())
}

func CheckBadRequestErrorResult(t *testing.T, result interface{}, errorMsg string) {
	require.NotNil(t, result)
	errorRes, isErrorResult := result.(ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	badReqErr, isBadReqErr := errorRes.Err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, errorMsg, badReqErr.Error())
}

func CheckUnexpectedErrorResult(t *testing.T, result interface{}, errorMsg string) {
	require.NotNil(t, result)
	errorRes, isErrorResult := result.(ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unexpErr, isUnexError := errorRes.Err.(*appErrors.UnexpectedError)
	require.Equal(t, true, isUnexError, "should be an unexpected error")
	assert.Equal(t, errorMsg, unexpErr.Error())
}

func CheckUnauthorizedErrorErrorResult(t *testing.T, result interface{}, errorMsg string) {
	require.NotNil(t, result)
	errorRes, isErrorResult := result.(ErrorResult)
	require.Equal(t, true, isErrorResult, "should be an error result")

	unauthpErr, isUnauthError := errorRes.Err.(*appErrors.UnauthorizedError)
	require.Equal(t, true, isUnauthError, "should be an unauthorized error")
	assert.Equal(t, errorMsg, unauthpErr.Error())
}
