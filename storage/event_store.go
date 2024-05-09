package storage

import (
	"context"

	"github.com/honestbank/event-driver/event"
)

// EventStore persists an event by key & source, and looks up an event by key+source, or a collection of events by key.
type EventStore interface {
	ListSourcesByKey(ctx context.Context, key string) ([]string, error)
	LookUp(ctx context.Context, key, source string) (*event.Message, error)
	LookUpByKey(ctx context.Context, key string) ([]*event.Message, error)
	Persist(ctx context.Context, key, source, content string) error
}
