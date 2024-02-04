package storage

import (
	"context"

	"github.com/lukecold/event-driver/event"
)

// GoogleCloudEventStore persists the contents in GCS, which requires consistent connections to Google Cloud.
type GoogleCloudEventStore struct {
}

func (g GoogleCloudEventStore) Persist(ctx context.Context, key, source, content string) error {
	//TODO implement me
	panic("implement me")
}

func (g GoogleCloudEventStore) LookUp(ctx context.Context, key, source string) (event.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoogleCloudEventStore) LookUpByKey(ctx context.Context, key string) ([]event.Message, error) {
	//TODO implement me
	panic("implement me")
}
