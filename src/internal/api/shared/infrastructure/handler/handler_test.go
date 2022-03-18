//go:build !e2e
// +build !e2e

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHandlerServeHTTP(t *testing.T) {
	t.Run("Returns 200 when no error", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.OkResult{Content: nil, StatusCode: http.StatusOK}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Returns 200 with content when no error", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			obj := struct {
				Field1 string
				Field2 string
			}{Field1: "a", Field2: "b"}
			return results.OkResult{Content: obj, StatusCode: http.StatusOK}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		want := "{\"Field1\":\"a\",\"Field2\":\"b\"}\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Returns 500 when an unexpected error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{Err: &appErrors.UnexpectedError{Msg: "error", InternalError: errors.New("msg")}}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "error\n", string(response.Body.String()))
	})

	t.Run("Returns 404 when a not found error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{Err: gorm.ErrRecordNotFound}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
		assert.Equal(t, "Not found\n", string(response.Body.String()))
	})

	t.Run("Returns 400 when a bad request error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{Err: &appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", "id")}}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
		assert.Equal(t, "\"id\" is not a valid id\n", string(response.Body.String()))
	})

	t.Run("Returns 401 when an unauthorized error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{Err: &appErrors.UnauthorizedError{Msg: "wadus"}}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "wadus\n", string(response.Body.String()))
	})

	t.Run("Returns 500 when an unhandled error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{Err: errors.New("wadus")}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "Internal error\n", string(response.Body.String()))
	})
}

func TestHandlerParseBody(t *testing.T) {
	handler := Handler{}

	t.Run("Returns a bad request error when the body is nil", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		res := handler.ParseBody(request, nil)

		assert.Error(t, res)
		badReqErr, isBadReqErr := res.(*appErrors.BadRequestError)
		require.Equal(t, true, isBadReqErr, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())
	})

	t.Run("Returns a bad request error when the body is invalid", func(t *testing.T) {
		body, _ := json.Marshal("")

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		res := handler.ParseBody(request, nil)

		assert.Error(t, res)
		badReqErr, isBadReqErr := res.(*appErrors.BadRequestError)
		require.Equal(t, true, isBadReqErr, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())
	})

	t.Run("Returns nil when the body is valid", func(t *testing.T) {
		body, _ := json.Marshal("text")

		data := ""
		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		res := handler.ParseBody(request, &data)

		assert.Nil(t, res)
		assert.Equal(t, "text", data)
	})
}
