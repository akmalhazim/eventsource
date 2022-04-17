package eventsource

type EventType string

type EventFactoryFn func() Event

type Event interface {
	EventType() EventType
}

var (
	events = make(map[EventType]EventFactoryFn)
)

func RegisterEvent(eventFactory EventFactoryFn) {
	event := eventFactory()
	if _, ok := events[event.EventType()]; ok {
		panic("Event is already registered")
	}

	events[event.EventType()] = eventFactory
}

func BuildEvent(eventType EventType) Event {
	if eventFactory, ok := events[eventType]; ok {
		return eventFactory()
	}

	panic("Event is not registered")
}
