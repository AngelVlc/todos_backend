package handlers

import (
	"os"
	"testing"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}

func CheckBadRequestErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	assert.Equal(t, true, isErrorResult, "should be an error result")

	badReqErr, isBadReqErr := errorRes.err.(*appErrors.BadRequestError)
	assert.Equal(t, true, isBadReqErr, "should be a bad request error")
	if isBadReqErr {
		assert.Equal(t, errorMsg, badReqErr.Error())
	}
}

func CheckUnexpectedErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	assert.Equal(t, true, isErrorResult, "should be an error result")

	unexpErr, isUnexError := errorRes.err.(*appErrors.UnexpectedError)
	assert.Equal(t, true, isUnexError, "should be an unexpected error")
	if isUnexError {
		assert.Equal(t, errorMsg, unexpErr.Error())
	}
}

func CheckUnauthorizedErrorErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	assert.Equal(t, true, isErrorResult, "should be an error result")

	unauthpErr, isUnauthError := errorRes.err.(*appErrors.UnauthorizedError)
	assert.Equal(t, true, isUnauthError, "should be an unauthorized error")
	if isUnauthError {
		assert.Equal(t, errorMsg, unauthpErr.Error())
	}
}
