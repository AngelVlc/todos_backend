package domain

import (
	"context"
)

type ConfigRepository interface {
	ExistsAllowedOrigin(ctx context.Context, origin Origin) (bool, error)
	GetAllAllowedOrigins(ctx context.Context) ([]AllowedOrigin, error)
	CreateAllowedOrigin(ctx context.Context, allowedOrigin *AllowedOrigin) error
	DeleteAllowedOrigin(ctx context.Context, id int32) error
}
