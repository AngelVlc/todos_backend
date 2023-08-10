package search

import "github.com/stretchr/testify/mock"

type MockedSearchIndexClient struct {
	mock.Mock
}

func NewMockedSearchIndexClient() *MockedSearchIndexClient {
	return &MockedSearchIndexClient{}
}

func (c *MockedSearchIndexClient) SaveObjects(objects interface{}) error {
	args := c.Called(objects)

	return args.Error(0)
}

func (m *MockedSearchIndexClient) DeleteObject(objectID string) error {
	args := m.Called(objectID)

	return args.Error(0)
}
