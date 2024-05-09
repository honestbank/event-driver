package transformer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers/transformer"
)

const (
	key     = "test-key"
	source  = "test-source"
	content = "test-content"
)

func TestEraseContentFromSources(t *testing.T) {
	eraseContentFromSources := transformer.EraseContentFromSources("source1", "source2")

	inputSourceToTransformedContent := map[string]string{
		"source1": "",
		"source2": "",
		"source3": content,
	}

	for inputSource, expectedTransformedContent := range inputSourceToTransformedContent {
		testName := fmt.Sprintf("transform %s", inputSource)
		t.Run(testName, func(t *testing.T) {
			transformed, err := eraseContentFromSources.Transform(event.NewMessage(key, inputSource, content))
			assert.NoError(t, err)
			assert.Equal(t, key, transformed.GetKey())
			assert.Equal(t, inputSource, transformed.GetSource())
			assert.Equal(t, expectedTransformedContent, transformed.GetContent())
		})
	}
}

func TestIdentityRule(t *testing.T) {
	identity := transformer.Identity()
	message := event.NewMessage(key, source, content)

	transformed, err := identity.Transform(message)
	assert.NoError(t, err)
	assert.Equal(t, message, transformed)
}

func TestRenameSourcesRule(t *testing.T) {
	aliasMap := map[string][]string{
		source: {"alias1", "alias2"},
	}
	renameSources, err := transformer.RenameSources(aliasMap)
	assert.NoError(t, err)

	inputSourceToTransformedSource := map[string]string{
		"alias1":    source,
		"alias2":    source,
		"not_alias": "not_alias",
	}

	for inputSource, expectedTransformedSource := range inputSourceToTransformedSource {
		testName := fmt.Sprintf("transform %s", inputSource)
		t.Run(testName, func(t *testing.T) {
			transformed, err := renameSources.Transform(event.NewMessage(key, inputSource, content))
			assert.NoError(t, err)
			assert.Equal(t, key, transformed.GetKey())
			assert.Equal(t, expectedTransformedSource, transformed.GetSource())
			assert.Equal(t, content, transformed.GetContent())
		})
	}
}
