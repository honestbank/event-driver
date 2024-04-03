package cache

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/storage"
)

// cache persists the input in a storage, and let user decide what to do in case of a cache hit.
type cache struct {
	conflictResolver ConflictResolver
	keyExtractor     KeyExtractor
	storage          storage.EventStore
}

func New(storage storage.EventStore) *cache {
	return &cache{
		conflictResolver: SkipOnConflict(),
		keyExtractor:     ExtractMessageKey(),
		storage:          storage,
	}
}

func (c *cache) WithKeyExtractor(keyExtractor KeyExtractor) *cache {
	c.keyExtractor = keyExtractor

	return c
}

func (c *cache) WithConflictResolver(conflictResolver ConflictResolver) *cache {
	c.conflictResolver = conflictResolver

	return c
}

func (c *cache) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	key, err := c.keyExtractor.Extract(in)
	if err != nil {
		return err
	}
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
