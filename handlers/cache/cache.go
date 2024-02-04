package cache

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/storage"
)

// ConflictResolver implements handlers.Handler that resolves the case
// when the input matches an existing record in the cache.
type ConflictResolver interface {
	Resolve(ctx context.Context, in event.Message, next handlers.CallNext) error
}

// cache persists the input in a storage, and let user decide what to do in case of a cache hit.
type cache struct {
	storage          storage.EventStore
	conflictResolver ConflictResolver
}

func New(storage storage.EventStore, conflictResolver ConflictResolver) *cache {
	return &cache{
		storage:          storage,
		conflictResolver: conflictResolver,
	}
}

func (c *cache) Process(ctx context.Context, in event.Message, next handlers.CallNext) error {
	key := in.GetKey()
	source := in.GetSource()
	message, err := c.storage.LookUp(ctx, key, source)
	if err != nil {
		return err
	}
	// cache hit
	if message != nil {
		return c.conflictResolver.Resolve(ctx, message, next)
	}

	// persist input message by key & source
	err = c.storage.Persist(ctx, key, source, in.GetContent())
	if err != nil {
		return err
	}

	return next.Call(ctx, in)
}
