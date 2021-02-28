package results

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type OkResult struct {
	Content    interface{}
	StatusCode int
}

func (r OkResult) IsError() bool {
	return false
}

func CheckOkResult(t *testing.T, result interface{}, expectedStatus int) *OkResult {
	require.NotNil(t, result)
	okRes, isOkResult := result.(OkResult)
	require.Equal(t, true, isOkResult, "should be an ok result")
	assert.Equal(t, expectedStatus, okRes.StatusCode)
	return &okRes
}
