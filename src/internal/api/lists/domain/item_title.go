package domain

import (
	"fmt"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type ItemTitle string

const title_max_length = 50

func NewItemTitle(title string) (ItemTitle, error) {
	if len(title) == 0 {
		return "", &appErrors.BadRequestError{Msg: "The item title can not be empty"}
	}

	if len(title) > title_max_length {
		return "", &appErrors.BadRequestError{Msg: fmt.Sprintf("The item title can not have more than %v characters", title_max_length)}
	}

	return ItemTitle(title), nil
}
