package cache_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/honestbank/event-driver/handlers/options"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/handlers/cache"
	"github.com/honestbank/event-driver/mocks"
	"github.com/honestbank/event-driver/storage"
)

func TestCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any())
		conflictResolver.EXPECT().Resolve(ctx, input1, callNext).Times(1) // trigger conflictResolver when cache hit

		logs := &strings.Builder{}
		handler := cache.New(eventStore, options.WithLogLevel(slog.LevelDebug), options.WithLogWriter(logs)).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		assert.Contains(t, logs.String(), "cache not hit")
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
		assert.Contains(t, logs.String(), "cache hit")
	})

	t.Run("cache not hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key1", "source", "content1")
		input2 := event.NewMessage("key2", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).Times(2)

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("failed to extract key", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key", "source", "content")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := mocks.NewMockKeyExtractor(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		keyExtractor.EXPECT().Extract(gomock.Any()).Return("", errors.New("test error"))

		logs := &strings.Builder{}
		handler := cache.New(eventStore, options.WithLogWriter(logs)).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input, callNext)
		assert.Error(t, err)
		assert.Contains(t, logs.String(), "failed to extract cache key")
		assert.Contains(t, logs.String(), `"error":"test error"`)
	})

	t.Run("failed to resolve", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		conflictResolver.EXPECT().Resolve(ctx, input1, callNext).Return(errors.New("test"))

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to pass to next", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key", "source", "content")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to lookup", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key1", "source", "content")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := mocks.NewMockEventStore(ctrl)

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		eventStore.EXPECT().LookUp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test")) // cache not hit on the first call

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to persist", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key", "source", "content")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.GetMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := mocks.NewMockEventStore(ctrl)

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		eventStore.EXPECT().LookUp(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil) // cache not hit on the first call
		eventStore.EXPECT().Persist(gomock.Any(), "key", "source", "content").Return(errors.New("test"))

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input, callNext)
		assert.Error(t, err)
	})
}
