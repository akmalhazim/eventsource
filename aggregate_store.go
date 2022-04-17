package eventsource

import (
	"context"

	"github.com/google/uuid"
)

type AggregateStore interface {
	Load(context.Context, AggregateType, uuid.UUID) (Aggregate, error)
	Save(context.Context, Aggregate) error
}

// store.Load("Transaction", "xxxx-xxxx-xxxx-xxxx")
