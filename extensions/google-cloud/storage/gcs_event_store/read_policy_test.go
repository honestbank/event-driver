package gcs_event_store_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/honestbank/event-driver/extensions/google-cloud/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/iterator"

	"github.com/honestbank/event-driver/extensions/google-cloud/storage/gcs_event_store"
)

func TestTakeFirstCreatedPolicy(t *testing.T) {
	t.Run("take the first created object", func(t *testing.T) {
		items := makeItems(5)
		ctrl := gomock.NewController(t)
		mockIterator := mocks.NewMockObjectIterator(ctrl)
		mockIterator.EXPECT().Next().DoAndReturn(func() (*gcs.ObjectAttrs, error) {
			if len(items) == 0 {
				return nil, iterator.Done
			}
			item := items[0]
			items = items[1:]

			return item, nil
		}).Times(len(items) + 1)

		policy := gcs_event_store.TakeFirstCreated()
		result, err := policy.Apply(mockIterator)
		assert.NoError(t, err)
		assert.Equal(t, "0", result.Name)
	})

	t.Run("fail in iteration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockIterator := mocks.NewMockObjectIterator(ctrl)
		mockIterator.EXPECT().Next().DoAndReturn(func() (*gcs.ObjectAttrs, error) {
			return nil, errors.New("test")
		})

		policy := gcs_event_store.TakeFirstCreated()
		_, err := policy.Apply(mockIterator)
		assert.Error(t, err)
	})

	t.Run("fail if iterator is  nil", func(t *testing.T) {
		policy := gcs_event_store.TakeFirstCreated()
		_, err := policy.Apply(nil)
		assert.Error(t, err)
	})
}

func TestTakeLastCreatedPolicy(t *testing.T) {
	t.Run("take the last created object", func(t *testing.T) {
		items := makeItems(5)
		ctrl := gomock.NewController(t)
		mockIterator := mocks.NewMockObjectIterator(ctrl)
		mockIterator.EXPECT().Next().DoAndReturn(func() (*gcs.ObjectAttrs, error) {
			if len(items) == 0 {
				return nil, iterator.Done
			}
			item := items[0]
			items = items[1:]

			return item, nil
		}).Times(len(items) + 1)

		policy := gcs_event_store.TakeLastCreated()
		result, err := policy.Apply(mockIterator)
		assert.NoError(t, err)
		assert.Equal(t, "4", result.Name)
	})

	t.Run("fail in iteration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockIterator := mocks.NewMockObjectIterator(ctrl)
		mockIterator.EXPECT().Next().DoAndReturn(func() (*gcs.ObjectAttrs, error) {
			return nil, errors.New("test")
		})

		policy := gcs_event_store.TakeLastCreated()
		_, err := policy.Apply(mockIterator)
		assert.Error(t, err)
	})

	t.Run("fail if iterator is  nil", func(t *testing.T) {
		policy := gcs_event_store.TakeLastCreated()
		_, err := policy.Apply(nil)
		assert.Error(t, err)
	})
}

func makeItems(number int) []*gcs.ObjectAttrs {
	timeBase := time.Now()
	objects := make([]*gcs.ObjectAttrs, 0, number)
	for i := 0; i < number; i++ {
		objects = append(objects, &gcs.ObjectAttrs{
			Name:    strconv.Itoa(i),
			Created: timeBase.Add(time.Duration(i) * time.Second),
		})
	}

	return objects
}
