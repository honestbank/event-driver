package cache_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/cache"
	"github.com/lukecold/event-driver/mocks"
	"github.com/lukecold/event-driver/storage"
)

type messageContent struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

type extractField1 struct{}

func (k *extractField1) Extract(in *event.Message) (string, error) {
	content := in.GetContent()
	var msg messageContent
	err := json.Unmarshal([]byte(content), &msg)
	if err != nil {
		return "", err
	}

	return msg.Field1, nil
}

func TestKeyExtractor(t *testing.T) {
	ctx := context.TODO()

	t.Run("default extractor", func(t *testing.T) {
		input1 := event.NewMessage("key", "source", "content1")
		input2 := event.NewMessage("key", "source", "content2")

		keyExtractor := cache.ExtractMessageKey()
		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		// callNext got triggered on input1, but skipped on input2
		callNext.EXPECT().Call(ctx, input1)

		handler := cache.New(eventStore).WithKeyExtractor(keyExtractor)

		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		val, err := eventStore.LookUp(context.TODO(), "key", "source")
		assert.NoError(t, err)
		assert.EqualValues(t, input1, val)

		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.EqualValues(t, input1, val)
	})

	t.Run("extractor json-typed content", func(t *testing.T) {
		input1 := event.NewMessage("key", "source", `{"field1":"val1","field2":"val2"}`)
		input2 := event.NewMessage("key", "source", `{"field1":"val1","field3":"val3"}`)

		keyExtractor := &extractField1{}
		ctrl := gomock.NewController(t)
		callNext := mocks.NewMockCallNext(ctrl)
		eventStore := storage.NewInMemoryStore()

		// callNext got triggered on input1, but skipped on input2
		callNext.EXPECT().Call(ctx, input1)

		handler := cache.New(eventStore).WithKeyExtractor(keyExtractor)

		err := handler.Process(ctx, input1, callNext)
		assert.NoError(t, err)
		val, err := eventStore.LookUp(context.TODO(), "val1", "source")
		assert.NoError(t, err)
		assert.EqualValues(t, event.NewMessage("val1", "source", `{"field1":"val1","field2":"val2"}`), val)

		err = handler.Process(ctx, input2, callNext)
		assert.NoError(t, err)
		val, err = eventStore.LookUp(context.TODO(), "val1", "source")
		assert.NoError(t, err)
		assert.EqualValues(t, event.NewMessage("val1", "source", `{"field1":"val1","field2":"val2"}`), val)
	})
}
