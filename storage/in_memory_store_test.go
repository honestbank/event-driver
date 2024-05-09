package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/storage"
)

const (
	key1    = "key1"
	key2    = "key2"
	source1 = "source1"
	source2 = "source2"
)

func TestNewInMemoryStore(t *testing.T) {
	inMemoryStore := storage.NewInMemoryStore()
	ctx := context.TODO()

	// persist
	err := inMemoryStore.Persist(ctx, key1, source1, "content1-1")
	assert.NoError(t, err)
	err = inMemoryStore.Persist(ctx, key1, source2, "content1-2")
	assert.NoError(t, err)
	err = inMemoryStore.Persist(ctx, key2, source1, "content2-1")
	assert.NoError(t, err)

	// list sources by key
	sources, err := inMemoryStore.ListSourcesByKey(ctx, key1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{source1, source2}, sources)
	sources, err = inMemoryStore.ListSourcesByKey(ctx, key2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{source1}, sources)
	sources, err = inMemoryStore.ListSourcesByKey(ctx, "unknown-key")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, sources)

	// look up contents by key
	contents, err := inMemoryStore.LookUpByKey(ctx, key1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*event.Message{
		event.NewMessage(key1, source1, "content1-1"),
		event.NewMessage(key1, source2, "content1-2"),
	}, contents)
	contents, err = inMemoryStore.LookUpByKey(ctx, key2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*event.Message{event.NewMessage(key2, source1, "content2-1")}, contents)
	contents, err = inMemoryStore.LookUpByKey(ctx, "unknown-key")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{}, contents)

	// look up contents
	content, err := inMemoryStore.LookUp(ctx, key1, source1)
	assert.NoError(t, err)
	assert.Equal(t, event.NewMessage(key1, source1, "content1-1"), content)
	content, err = inMemoryStore.LookUp(ctx, key1, source2)
	assert.NoError(t, err)
	assert.Equal(t, event.NewMessage(key1, source2, "content1-2"), content)
	content, err = inMemoryStore.LookUp(ctx, key2, source1)
	assert.NoError(t, err)
	assert.Equal(t, event.NewMessage(key2, source1, "content2-1"), content)
	content, err = inMemoryStore.LookUp(ctx, key2, source2)
	assert.NoError(t, err)
	assert.Nil(t, content)
}
