package aggregatestore

import (
	"context"

	"github.com/akmalhazim/eventsource"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAggregateStore struct {
	collection *mongo.Collection
}

func (store *mongoAggregateStore) Load(ctx context.Context, aggregateType eventsource.AggregateType, aggregateID uuid.UUID) (eventsource.Aggregate, error) {
	cursor, err := store.collection.Find(ctx, bson.M{
		"aggregateType": aggregateType,
		"aggregateId":   aggregateID,
	})
	if err != nil {
		// TODO check if DocumentNotExist exception, then we factory create a new empty aggregate
		return nil, err
	}
	defer cursor.Close(ctx)

	aggregate := eventsource.BuildAggregate(aggregateType, aggregateID)
	for cursor.Next(ctx) {
		serializedEvent := new(eventsource.SerializedEvent)
		err := cursor.Decode(serializedEvent)
		if err != nil {
			return nil, err
		}
		event := eventsource.BuildEvent(serializedEvent.EventType)
		err = bson.Unmarshal(serializedEvent.Data, event)
		if err != nil {
			return nil, err
		}

		// WIP apply all events to the built Aggregate
		aggregate.HandleEvent(event)
	}

	return aggregate, nil
}

func (store *mongoAggregateStore) Save(ctx context.Context, aggregate eventsource.Aggregate) error {
	events := aggregate.UncommittedEvents()

	for _, event := range events {
		// TODO implement database transactions
		_, err := store.collection.InsertOne(ctx, bson.M{
			"aggregateType": aggregate.AggregateType(),
			"aggregateId":   aggregate.AggregateID(),
			"eventType":     event.EventType(),
			"data":          event,
		})
		if err != nil {
			return err
		}
	}

	for _, event := range events {
		aggregate.HandleEvent(event)
	}

	return nil
}

func NewMongoAggregateStore(db *mongo.Database) eventsource.AggregateStore {
	return &mongoAggregateStore{
		collection: db.Collection("events"),
	}
}
