package cache_test

import (
	"context"
	"testing"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/cache"
	"github.com/lukecold/event-driver/mocks"
	"github.com/lukecold/event-driver/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSkipOnConflict(t *testing.T) {
	ctx := context.TODO()
	input1 := event.NewMessage("key", "source", "content1")
	input2 := event.NewMessage("key", "source", "content2")

	skipOnConflict := cache.SkipOnConflict()
	ctrl := gomock.NewController(t)
	callNext := mocks.NewMockCallNext(ctrl)
	eventStore := storage.NewInMemoryStore()

	// callNext got triggered on input1, but skipped on input2
	callNext.EXPECT().Call(ctx, input1)

	handler := cache.New(eventStore, skipOnConflict)
	err := handler.Process(ctx, input1, callNext)
	assert.NoError(t, err)
	err = handler.Process(ctx, input2, callNext)
	assert.NoError(t, err)
}
