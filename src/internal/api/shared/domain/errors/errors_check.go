package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CheckErrorMsg(t *testing.T, err error, errorMsg string) {
	if assert.Error(t, err) {
		assert.Equal(t, errorMsg, err.Error())
	}
}

func CheckUnexpectedError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	require.NotNil(t, err)
	unexpectErr, isUnexpectError := err.(*UnexpectedError)
	require.True(t, isUnexpectError, "should be an unexpected error")
	CheckErrorMsg(t, unexpectErr, errorMsg)

	if len(internalErrorMsg) > 0 {
		CheckErrorMsg(t, unexpectErr.InternalError, internalErrorMsg)
	}
}

func CheckUnathorizedError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	require.NotNil(t, err)
	unauthErr, isUnauthErr := err.(*UnauthorizedError)
	require.True(t, isUnauthErr, "should be an unauthorized error")
	CheckErrorMsg(t, unauthErr, errorMsg)

	if len(internalErrorMsg) > 0 {
		CheckErrorMsg(t, unauthErr.InternalError, internalErrorMsg)
	}
}

func CheckBadRequestError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	require.NotNil(t, err)
	badReqErr, isBadReqErr := err.(*BadRequestError)
	require.True(t, isBadReqErr, "should be a bad request error")
	CheckErrorMsg(t, badReqErr, errorMsg)

	if len(internalErrorMsg) > 0 {
		CheckErrorMsg(t, badReqErr.InternalError, internalErrorMsg)
	}
}
