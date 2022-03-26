package domain

import (
	"context"

	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
)

type Origin string

func NewUserName(origin string) (Origin, error) {
	if len(origin) == 0 {
		return "", &appErrors.BadRequestError{Msg: "Origin can not be empty"}
	}

	return Origin(origin), nil
}

func (o Origin) CheckIfAlreadyExists(ctx context.Context, repo ConfigRepository) error {
	exists, err := repo.ExistsAllowedOrigin(ctx, o)
	if err != nil {
		return err
	}

	if exists {
		return &appErrors.BadRequestError{Msg: "The origin already exists", InternalError: nil}
	}

	return nil
}
