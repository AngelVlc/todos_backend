package events

type Subscriber interface {
	Subscribe()
	Start()
}
