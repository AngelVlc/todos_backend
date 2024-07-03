package domain

import "context"

type CategoriesRepository interface {
	FindCategory(ctx context.Context, query CategoryRecord) (*CategoryRecord, error)
	ExistsCategory(ctx context.Context, query CategoryRecord) (bool, error)
	GetCategories(ctx context.Context, query CategoryRecord) (CategoryRecords, error)
	CreateCategory(ctx context.Context, record *CategoryRecord) error
	DeleteCategory(ctx context.Context, query CategoryRecord) error
	UpdateCategory(ctx context.Context, record *CategoryRecord) error
}
