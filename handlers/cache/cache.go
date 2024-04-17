package cache

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/log"
	"github.com/lukecold/event-driver/storage"
)

// cache persists the input in a storage, and let user decide what to do in case of a cache hit.
type cache struct {
	cacheKeyExtractor KeyExtractor
	conflictResolver  ConflictResolver
	logger            *log.Logger
	storage           storage.EventStore
}

func New(storage storage.EventStore) *cache {
	return &cache{
		cacheKeyExtractor: GetMessageKey(),
		conflictResolver:  SkipOnConflict(),
		logger:            log.New("cache"),
		storage:           storage,
	}
}

func (c *cache) WithKeyExtractor(keyExtractor KeyExtractor) *cache {
	c.cacheKeyExtractor = keyExtractor

	return c
}

func (c *cache) WithConflictResolver(conflictResolver ConflictResolver) *cache {
	c.conflictResolver = conflictResolver

	return c
}

func (c *cache) Verbose() *cache {
	c.logger.Verbose()

	return c
}

func (c *cache) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	key, err := c.cacheKeyExtractor.Extract(in)
	if err != nil {
		c.logger.Error("failed to extract cache key for message_key=%s", in.GetKey())

		return err
	}
	source := in.GetSource()
	message, err := c.storage.LookUp(ctx, key, source)
	if err != nil {
		c.logger.Error("failed to look up message with cache_key=%s", key)

		return err
	}
	// cache hit
	if message != nil {
		c.logger.Info("cache hit on cache_key=%s", key)

		return c.conflictResolver.Resolve(ctx, message, next)
	}

	// persist input message by key & source
	err = c.storage.Persist(ctx, key, source, in.GetContent())
	if err != nil {
		c.logger.Error("failed to persist message with cache_key=%s, source=%s", key, in.GetSource())

		return err
	}
	c.logger.Debug("cache not hit on cache_key=%s", key)

	return next.Call(ctx, in)
}
