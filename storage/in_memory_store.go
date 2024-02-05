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

func (i *InMemoryStore) Persist(ctx context.Context, key, source, content string) error {
	if _, isKeyExist := (*i)[key]; !isKeyExist {
		(*i)[key] = make(map[string]string)
	}
	(*i)[key][source] = content

	return nil
}

func (i *InMemoryStore) LookUp(ctx context.Context, key, source string) (*event.Message, error) {
	content, isHit := (*i)[key][source]
	if !isHit {
		return nil, nil
	}

	return event.NewMessage(key, source, content), nil
}

func (i *InMemoryStore) LookUpByKey(ctx context.Context, key string) ([]*event.Message, error) {
	results := (*i)[key]
	messages := make([]*event.Message, 0, len(results))
	for source, content := range results {
		messages = append(messages, event.NewMessage(key, source, content))
	}

	return messages, nil
}
