package cache_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/cache"
	"github.com/lukecold/event-driver/mocks"
	"github.com/lukecold/event-driver/storage"
)

func TestCache(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.ExtractMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()
		conflictResolver.EXPECT().Resolve(ctx, input1, callNext).Times(1) // trigger conflictResolver when cache hit

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("cache not hit", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key1", "source", "content1")
		input2 := event.NewMessage("key2", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.ExtractMessageKey()
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		callNext.EXPECT().Call(gomock.Any(), gomock.Any()).AnyTimes()

		handler := cache.New(eventStore).
			WithConflictResolver(conflictResolver).
			WithKeyExtractor(keyExtractor)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("failed to resolve", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		ctrl := gomock.NewController(t)
		conflictResolver := mocks.NewMockConflictResolver(ctrl)
		keyExtractor := cache.ExtractMessageKey()
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
		keyExtractor := cache.ExtractMessageKey()
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
		keyExtractor := cache.ExtractMessageKey()
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
		keyExtractor := cache.ExtractMessageKey()
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
