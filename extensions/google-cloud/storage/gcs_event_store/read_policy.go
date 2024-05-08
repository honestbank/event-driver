package gcs_event_store

import (
	"errors"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type ObjectIterator interface {
	Next() (*gcs.ObjectAttrs, error)
}

// ReadPolicy defines the actions that the GCSEventStore takes when it got the object iterator via list query.
type ReadPolicy interface {
	Apply(ObjectIterator) (*gcs.ObjectAttrs, error)
}

// takeFirstCreated take the earliest created object under the same key/source/ path
type takeFirstCreated struct{}

func TakeFirstCreated() ReadPolicy {
	return takeFirstCreated{}
}

func (t takeFirstCreated) Apply(objectIterator ObjectIterator) (*gcs.ObjectAttrs, error) {
	if objectIterator == nil {
		return nil, errors.New("objectIterator is nil")
	}

	var result *gcs.ObjectAttrs
	for {
		object, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		if result == nil || object.Created.Before(result.Created) {
			result = object
		}
	}

	return result, nil
}

// takeLastCreated take the latest created object under the same key/source/ path
type takeLastCreated struct{}

func TakeLastCreated() ReadPolicy {
	return takeLastCreated{}
}

func (t takeLastCreated) Apply(objectIterator ObjectIterator) (*gcs.ObjectAttrs, error) {
	if objectIterator == nil {
		return nil, errors.New("objectIterator is nil")
	}

	var result *gcs.ObjectAttrs
	for {
		object, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		if result == nil || object.Created.After(result.Created) {
			result = object
		}
	}

	return result, nil
}
