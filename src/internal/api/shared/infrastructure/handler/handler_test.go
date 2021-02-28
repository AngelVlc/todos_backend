//+build !e2e

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Returns 200 when no error", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.OkResult{nil, http.StatusOK}
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
			return results.OkResult{obj, http.StatusOK}
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
			return results.ErrorResult{&appErrors.UnexpectedError{Msg: "error", InternalError: errors.New("msg")}}
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
			return results.ErrorResult{&appErrors.NotFoundError{Model: "model"}}
		}

		handler := Handler{
			HandlerFunc: f,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
		assert.Equal(t, "model not found\n", string(response.Body.String()))
	})

	t.Run("Returns 400 when a bad request error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
			return results.ErrorResult{&appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", "id")}}
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
			return results.ErrorResult{&appErrors.UnauthorizedError{Msg: "wadus"}}
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
			return results.ErrorResult{errors.New("wadus")}
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
