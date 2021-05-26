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
	m.Called(eventName, data)
	m.mu.Unlock()
}

func (m *MockedEventBus) Subscribe(eventName string, ch DataChannel) {
	defer m.Wg.Done()

	m.mu.Lock()
	m.Called(eventName, ch)
	m.mu.Unlock()
}
