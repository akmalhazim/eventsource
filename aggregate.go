package eventsource

import (
	"sync"

	"github.com/google/uuid"
)

type AggregateType string

type AggregateFactoryFn func(uuid.UUID) Aggregate

type Aggregate interface {
	AggregateType() AggregateType
	AggregateID() uuid.UUID
	HandleEvent(Event)
	UncommittedEvents() []Event
}

// Aggregate is thread-safe
type BaseAggregate struct {
	Aggregate
	sync.Mutex
	events []Event
}

func (aggregate *BaseAggregate) AppendEvent(event Event) {
	aggregate.Mutex.Lock()
	defer aggregate.Unlock()

	aggregate.events = append(aggregate.events, event)
}

func (aggregate *BaseAggregate) UncommittedEvents() []Event {
	return aggregate.events
}

// TODO add another abstract function for deriving Aggregate to HandleEvent

var aggregates = make(map[AggregateType]AggregateFactoryFn)

// RegisterAggregate should only run during initialization. It's not thread-safe.
func RegisterAggregate(aggregateFactory AggregateFactoryFn) {
	aggregate := aggregateFactory(uuid.New())

	if _, ok := aggregates[aggregate.AggregateType()]; ok {
		panic("Aggregate is already registered")
	}

	aggregates[aggregate.AggregateType()] = aggregateFactory
}

// BuildAggregate should run during initalization. This is important to build Aggregate(s) from scratch - our Mongo repository
func BuildAggregate(aggregateType AggregateType, aggregateID uuid.UUID) Aggregate {
	if aggregateFactory, ok := aggregates[aggregateType]; ok {
		return aggregateFactory(aggregateID)
	}

	panic("Aggregate is not registered")
}
