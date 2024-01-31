package transformer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/transformer"
	"github.com/lukecold/event-driver/mocks"
)

func TestTransformer(t *testing.T) {
	aliasMap := map[string][]string{
		source: {"alias1", "alias2"},
	}
	renameSources, err := transformer.RenameSources(aliasMap)
	assert.NoError(t, err)
	eraseContentFromSource := transformer.EraseContentFromSources(source)
	eraseContentFromSource1And2 := transformer.EraseContentFromSources("source1", "source2")
	eventMapper := transformer.New(renameSources, eraseContentFromSource).
		WithRules(eraseContentFromSource1And2)

	inputSourceToTransformedEvent := map[string]event.Message{
		"alias1":  event.NewMessage(key, source, ""),
		"alias2":  event.NewMessage(key, source, ""),
		"source1": event.NewMessage(key, "source1", ""),
		"source2": event.NewMessage(key, "source2", ""),
		"source3": event.NewMessage(key, "source3", content),
	}

	for inputSource, expectedTransformedEvent := range inputSourceToTransformedEvent {
		ctx := context.TODO()
		ctrl := gomock.NewController(t)
		next := mocks.NewMockCallNext(ctrl)
		next.EXPECT().Call(ctx, expectedTransformedEvent)

		message := event.NewMessage(key, inputSource, content)
		err := eventMapper.Process(ctx, message, next)
		assert.NoError(t, err)
	}
}
