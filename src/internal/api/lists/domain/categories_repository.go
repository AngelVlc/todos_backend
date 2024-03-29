package domain

import "context"

type CategoriesRepository interface {
	FindCategory(ctx context.Context, query CategoryEntity) (*CategoryEntity, error)
	ExistsCategory(ctx context.Context, query CategoryEntity) (bool, error)
	GetAllCategoriesForUser(ctx context.Context, userID int32) ([]*CategoryEntity, error)
	CreateCategory(ctx context.Context, list *CategoryEntity) (*CategoryEntity, error)
	DeleteCategory(ctx context.Context, query CategoryEntity) error
	UpdateCategory(ctx context.Context, list *CategoryEntity) (*CategoryEntity, error)
}
