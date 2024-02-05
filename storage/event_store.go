package storage

import (
	"context"

	"github.com/lukecold/event-driver/event"
)

// EventStore persists an event by key & source, and looks up an event by key+source, or a collection of events by key.
type EventStore interface {
	Persist(ctx context.Context, key, source, content string) error
	LookUp(ctx context.Context, key, source string) (*event.Message, error)
	LookUpByKey(ctx context.Context, key string) ([]*event.Message, error)
}
