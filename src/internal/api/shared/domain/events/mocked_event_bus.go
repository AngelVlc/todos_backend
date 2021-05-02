package events

import "github.com/stretchr/testify/mock"

type MockedEventBus struct {
	mock.Mock
}

func NewMockedEventBus() *MockedEventBus {
	return &MockedEventBus{}
}

func (m *MockedEventBus) Publish(topic string, data interface{}) {
	m.Called(topic, data)
}

func (m *MockedEventBus) Subscribe(topic string, ch DataChannel) {
	m.Called(topic, ch)
}
