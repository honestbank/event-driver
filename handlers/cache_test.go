package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/KeluDiao/event-driver/event"
	"github.com/KeluDiao/event-driver/handlers"
	"github.com/KeluDiao/event-driver/mocks"
	"github.com/KeluDiao/event-driver/storage"
)

func TestCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		cacheHitResolver.EXPECT().Resolve(ctx, input1, callNext).Times(1) // trigger cacheHitResolver when cache hit

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = cache.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("cache not hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key1", "source", "content1")
		input2 := event.NewMessage("key2", "source", "content2")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = cache.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("failed to resolve", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		cacheHitResolver.EXPECT().Resolve(ctx, input1, callNext).Return(errors.New("test"))

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = cache.Process(ctx, input2, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to pass to next", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key", "source", "content")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to lookup", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key1", "source", "content")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := mocks.NewMockEventStore(ctrl)

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		eventStore.EXPECT().LookUp(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")) // cache not hit on the first call

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to persist", func(t *testing.T) {
		ctx := context.TODO()
		input := event.NewMessage("key", "source", "content")

		ctrl := gomock.NewController(t)
		cacheHitResolver := mocks.NewMockCacheHitResolver(ctrl)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := mocks.NewMockEventStore(ctrl)

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		eventStore.EXPECT().LookUp(gomock.Any(), gomock.Any()).Return(nil, nil) // cache not hit on the first call
		eventStore.EXPECT().Persist("key", "source", "content").Return(errors.New("test"))

		cache := handlers.NewCache(eventStore, cacheHitResolver)
		err := cache.Process(ctx, input, callNext)
		assert.Error(t, err)
	})
}
