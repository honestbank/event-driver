package handlers

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/storage"
)

// CacheHitResolver resolves the case when the input matches a record in the cache
type CacheHitResolver interface {
	Resolve(ctx context.Context, in event.Message, next CallNext) error
}

// Cache persists the input in a storage, and let user decide what to do in case of a cache hit.
type Cache struct {
	storage          storage.EventStore
	cacheHitResolver CacheHitResolver
}

func NewCache(storage storage.EventStore, resolveCacheHit CacheHitResolver) *Cache {
	return &Cache{
		storage:          storage,
		cacheHitResolver: resolveCacheHit,
	}
}

func (c *Cache) Process(ctx context.Context, in event.Message, next CallNext) error {
	key := in.GetKey()
	source := in.GetSource()
	message, err := c.storage.LookUp(key, source)
	if err != nil {
		return err
	}
	// cache hit
	if message != nil {
		return c.cacheHitResolver.Resolve(ctx, message, next)
	}

	// persist input message by key & source
	err = c.storage.Persist(key, source, in.GetContent())
	if err != nil {
		return err
	}

	return next.Call(ctx, in)
}
