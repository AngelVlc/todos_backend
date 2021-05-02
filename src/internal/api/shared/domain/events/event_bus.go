package events

type EventBus interface {
	Publish(topic string, data interface{})
	Subscribe(topic string, ch DataChannel)
}
