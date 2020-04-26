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
	assert.Equal(t, errorMsg, badReqErr.Error())
}

func CheckUnexpectedErrorResult(t *testing.T, result interface{}, errorMsg string) {
	assert.NotNil(t, result)
	errorRes, isErrorResult := result.(errorResult)
	assert.Equal(t, true, isErrorResult, "should be an error result")

	unexpErr, isUnexError := errorRes.err.(*appErrors.UnexpectedError)
	assert.Equal(t, true, isUnexError, "should be a bad request error")
	assert.Equal(t, errorMsg, unexpErr.Error())
}
