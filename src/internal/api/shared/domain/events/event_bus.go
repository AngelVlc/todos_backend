package events

type EventBus interface {
	Publish(eventName string, data interface{})
	Subscribe(eventName string, ch DataChannel)
}
