package events

import (
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockedEventBus struct {
	mock.Mock
	mu sync.Mutex
	Wg sync.WaitGroup
}

func NewMockedEventBus() *MockedEventBus {
	return &MockedEventBus{}
}

func (m *MockedEventBus) Publish(eventName string, data interface{}) {
	defer m.Wg.Done()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Called(eventName, data)
}

func (m *MockedEventBus) Subscribe(eventName string, ch DataChannel) {
	defer m.Wg.Done()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Called(eventName, ch)
}
