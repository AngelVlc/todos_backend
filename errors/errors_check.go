package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckErrorMsg(t *testing.T, err error, errorMsg string) {
	if assert.Error(t, err) {
		assert.Equal(t, errorMsg, err.Error())
	}
}

func CheckUnexpectedError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	assert.NotNil(t, err)
	unexpectErr, isUnexpectError := err.(*UnexpectedError)
	assert.Equal(t, true, isUnexpectError, "should be an unexpected error")
	if isUnexpectError {
		CheckErrorMsg(t, unexpectErr, errorMsg)
		if len(internalErrorMsg) > 0 {
			CheckErrorMsg(t, unexpectErr.InternalError, internalErrorMsg)
		}
	}
}

func CheckUnathorizedError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	assert.NotNil(t, err)
	unauthErr, isUnauthErr := err.(*UnauthorizedError)
	assert.Equal(t, true, isUnauthErr, "should be an unauthorized error")
	if isUnauthErr {
		CheckErrorMsg(t, unauthErr, errorMsg)
		if len(internalErrorMsg) > 0 {
			CheckErrorMsg(t, unauthErr.InternalError, internalErrorMsg)
		}
	}
}

func CheckBadRequestError(t *testing.T, err interface{}, errorMsg string, internalErrorMsg string) {
	assert.NotNil(t, err)
	badReqErr, isBadReqErr := err.(*BadRequestError)
	assert.Equal(t, true, isBadReqErr, "should be a bad request error")
	if isBadReqErr {
		CheckErrorMsg(t, badReqErr, errorMsg)
		if len(internalErrorMsg) > 0 {
			CheckErrorMsg(t, badReqErr.InternalError, internalErrorMsg)
		}
	}
}
