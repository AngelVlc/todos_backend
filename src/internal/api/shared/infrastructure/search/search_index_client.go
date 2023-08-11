package search

type SearchIndexClient interface {
	SaveObjects(objects interface{}) error
	DeleteObject(objectID string) error
}
