package helpers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

func ParseBody(r *http.Request, dto interface{}) error {
	if r.Body == nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: nil}
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return nil
}
