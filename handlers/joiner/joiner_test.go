package joiner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/joiner"
	"github.com/lukecold/event-driver/mocks"
	"github.com/lukecold/event-driver/storage"
)

func TestJoiner(t *testing.T) {
	t.Run("condition met", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source1", "content1")
		input2 := event.NewMessage("key", "source2", "content2")

		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		condition := joiner.MatchAll("source1").
			And(joiner.MatchAny("source2", "source3"))
		eventStore := storage.NewInMemoryStore()

		expectedMessage := event.NewMessage("key", "composed-event", `{"source1":"content1","source2":"content2"}`)
		callNext.EXPECT().Call(gomock.Any(), expectedMessage).AnyTimes()

		handler := joiner.New(condition, eventStore)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("condition not met", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source2", "content2")
		input2 := event.NewMessage("key", "source3", "content3")

		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		condition := joiner.MatchAll("source1").
			And(joiner.MatchAny("source2", "source3"))
		eventStore := storage.NewInMemoryStore()

		handler := joiner.New(condition, eventStore)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
	})

	t.Run("failed to pass to next", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source1", "content1")
		input2 := event.NewMessage("key", "source2", "content2")

		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		condition := joiner.MatchAll("source1").
			And(joiner.MatchAny("source2", "source3"))
		eventStore := storage.NewInMemoryStore()

		expectedMessage := event.NewMessage("key", "composed-event", `{"source1":"content1","source2":"content2"}`)
		callNext.EXPECT().Call(gomock.Any(), expectedMessage).AnyTimes().Return(errors.New("test"))

		handler := joiner.New(condition, eventStore)
		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to persist", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source1", "content1")
		input2 := event.NewMessage("key", "source2", "content2")

		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		condition := joiner.MatchAll("source1").
			And(joiner.MatchAny("source2", "source3"))
		eventStore := mocks.NewMockEventStore(ctrl)

		eventStore.EXPECT().Persist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("test")).Times(2)

		handler := joiner.New(condition, eventStore)
		err := handler.Process(ctx, input1, callNext)
		assert.Error(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.Error(t, err)
	})

	t.Run("failed to lookup by key", func(t *testing.T) {
		ctx := context.TODO()
		input1 := event.NewMessage("key", "source1", "content1")
		input2 := event.NewMessage("key", "source2", "content2")

		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		condition := joiner.MatchAll("source1").
			And(joiner.MatchAny("source2", "source3"))
		eventStore := mocks.NewMockEventStore(ctrl)

		eventStore.EXPECT().Persist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		eventStore.EXPECT().LookUpByKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(2)

		handler := joiner.New(condition, eventStore)
		err := handler.Process(ctx, input1, callNext)
		assert.Error(t, err)
		err = handler.Process(ctx, input2, callNext)
		assert.Error(t, err)
	})
}
