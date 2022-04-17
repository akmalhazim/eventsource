package eventsource

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type SerializedEvent struct {
	AggregateType string    `bson:"aggregateType"`
	AggregateID   uuid.UUID `bson:"aggregateId"`
	EventType     EventType `bson:"eventType"`
	Data          bson.Raw  `bson:"data"`
}
