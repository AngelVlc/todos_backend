package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestTokenHandler(t *testing.T) {
	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))

		result := TokenHandler(request, Handler{})

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body does not have user name", func(t *testing.T) {
		login := struct {
			Password string
		}{
			"pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(request, Handler{})

		CheckBadRequestErrorResult(t, result, "UserName is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body does not have password", func(t *testing.T) {
		login := struct {
			UserName string
		}{
			"wadus",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(request, Handler{})

		CheckBadRequestErrorResult(t, result, "Password is mandatory")
	})

	// t.Run("Should return an error result with an unexpexted error if getting the user fails", func(t *testing.T) {
	// 	login := models.Login{
	// 		UserName: "wadus",
	// 		Password: "pass",
	// 	}
	// 	body, _ := json.Marshal(login)

	// 	request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

	// 	result := TokenHandler(request, nil)

	// 	CheckUnexpectedErrorResult(t, result, "Password is mandatory")
	// })
}
