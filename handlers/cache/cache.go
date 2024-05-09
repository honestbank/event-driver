package cache

import (
	"context"
	"log/slog"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/handlers"
	"github.com/honestbank/event-driver/handlers/options"
	"github.com/honestbank/event-driver/storage"
)

// cache persists the input in a storage, and let user decide what to do in case of a cache hit.
type cache struct {
	cacheKeyExtractor KeyExtractor
	conflictResolver  ConflictResolver
	logger            *slog.Logger
	storage           storage.EventStore
}

func New(storage storage.EventStore, opts ...options.Option) *cache {
	cfg := options.DefaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &cache{
		cacheKeyExtractor: GetMessageKey(),
		conflictResolver:  SkipOnConflict(),
		logger: slog.New(slog.NewJSONHandler(cfg.GetLogWriter(), &slog.HandlerOptions{Level: cfg.GetLogLevel()})).
			With(slog.String("handler", "cache")),
		storage: storage,
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

func (c *cache) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	logger := c.logger.With(slog.String("key", in.GetKey()), slog.String("source", in.GetSource()))
	key, err := c.cacheKeyExtractor.Extract(in)
	if err != nil {
		logger.Error("failed to extract cache key", slog.Any("error", err))

		return err
	}
	logger = logger.With(slog.String("cache_key", key))
	source := in.GetSource()
	message, err := c.storage.LookUp(ctx, key, source)
	if err != nil {
		logger.Error("failed to look up message", slog.Any("error", err))

		return err
	}
	// cache hit
	if message != nil {
		logger.Info("cache hit")

		return c.conflictResolver.Resolve(ctx, message, next)
	}

	// persist input message by key & source
	err = c.storage.Persist(ctx, key, source, in.GetContent())
	if err != nil {
		logger.Error("failed to persist message", slog.Any("error", err))

		return err
	}
	logger.Debug("cache not hit")

	return next.Call(ctx, in)
}
