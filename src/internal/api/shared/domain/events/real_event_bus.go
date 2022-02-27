package events

import "sync"

// RealEventBus stores the information about subscribers interested for a particular event
type RealEventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func NewRealEventBus(subscribers map[string]DataChannelSlice) *RealEventBus {
	return &RealEventBus{
		subscribers: subscribers,
	}
}

func (eb *RealEventBus) Publish(eventName string, data interface{}) {
	eb.rm.RLock()

	if chans, found := eb.subscribers[eventName]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		// special thanks for /u/freesid who pointed it out
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: eventName}, channels)
	}

	eb.rm.RUnlock()
}

func (eb *RealEventBus) Subscribe(eventName string, ch DataChannel) {
	eb.rm.Lock()
	defer eb.rm.Unlock()

	if prev, found := eb.subscribers[eventName]; found {
		eb.subscribers[eventName] = append(prev, ch)
	} else {
		eb.subscribers[eventName] = append([]DataChannel{}, ch)
	}
}
