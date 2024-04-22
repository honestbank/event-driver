package storage

import (
	"context"

	"github.com/lukecold/event-driver/event"
)

type InMemoryStore map[string]map[string]string // key -> source -> content

func NewInMemoryStore() *InMemoryStore {
	internalMap := make(map[string]map[string]string)
	inMemoryStore := InMemoryStore(internalMap)

	return &inMemoryStore
}

func (i *InMemoryStore) ListSourcesByKey(_ context.Context, key string) ([]string, error) {
	results := (*i)[key]
	sources := make([]string, 0, len(results))
	for source := range results {
		sources = append(sources, source)
	}

	return sources, nil
}

func (i *InMemoryStore) LookUp(_ context.Context, key, source string) (*event.Message, error) {
	content, isHit := (*i)[key][source]
	if !isHit {
		return nil, nil
	}

	return event.NewMessage(key, source, content), nil
}

func (i *InMemoryStore) LookUpByKey(_ context.Context, key string) ([]*event.Message, error) {
	results := (*i)[key]
	messages := make([]*event.Message, 0, len(results))
	for source, content := range results {
		messages = append(messages, event.NewMessage(key, source, content))
	}

	return messages, nil
}

func (i *InMemoryStore) Persist(_ context.Context, key, source, content string) error {
	if _, isKeyExist := (*i)[key]; !isKeyExist {
		(*i)[key] = make(map[string]string)
	}
	(*i)[key][source] = content

	return nil
}
