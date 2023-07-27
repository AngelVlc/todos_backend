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
	"strings"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHandlerServeHTTP(t *testing.T) {
	t.Run("Returns a 400 with invalid body if the request requires an input but the body is empty", func(t *testing.T) {
		handler := Handler{
			RequestInput: &domain.CreateListInput{},
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
		assert.Equal(t, "Invalid body\n", string(response.Body.String()))
	})

	t.Run("Returns a 400 with invalid body if the request requires an input but the body is not the expected input", func(t *testing.T) {
		handler := Handler{
			RequestInput: &domain.CreateListInput{},
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
		assert.Equal(t, "Invalid body\n", string(response.Body.String()))
	})

	t.Run("Returns 200 when no error when the request requires an input and the body has the expected input", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.OkResult{Content: nil, StatusCode: http.StatusOK}
		}

		handler := Handler{
			HandlerFunc:  f,
			RequestInput: &domain.CreateListInput{},
		}

		listName, _ := domain.NewListNameValueObject("list1")
		createReq := domain.CreateListInput{Name: listName}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
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
