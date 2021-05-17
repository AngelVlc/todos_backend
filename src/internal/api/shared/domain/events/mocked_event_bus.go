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

func (m *MockedEventBus) Publish(topic string, data interface{}) {
	defer m.Wg.Done()

	m.mu.Lock()
	m.Called(topic, data)
	m.mu.Unlock()
}

func (m *MockedEventBus) Subscribe(topic string, ch DataChannel) {
	defer m.Wg.Done()

	m.mu.Lock()
	m.Called(topic, ch)
	m.mu.Unlock()
}
