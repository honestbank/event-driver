//go:build integration_test

package gcs_event_store_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/extensions/google-cloud/storage/gcs_event_store"
	"github.com/honestbank/event-driver/utils/compression"
)

const (
	key     = "key"
	source1 = "source1"
	source2 = "source2"
	content = "content"
)

func TestGCSEventStore(t *testing.T) {
	_ = os.Setenv("STORAGE_EMULATOR_HOST", "localhost:19081")
	folderName := "folder-name"

	t.Run("query empty bucket", func(t *testing.T) {
		bucket := "empty-bucket"
		setup(t, bucket)
		config := gcs_event_store.Config(bucket).WithFolder(folderName)
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, sources)

		message, err := eventStore.LookUp(context.TODO(), key, source1)
		assert.NoError(t, err)
		assert.Nil(t, message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []*event.Message{}, messageArray)
	})

	t.Run("with folder prefix", func(t *testing.T) {
		bucket := "with-prefix"
		setup(t, bucket)
		config := gcs_event_store.Config(bucket).WithFolder(folderName)
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source2, content)
		assert.NoError(t, err)
		err = eventStore.Persist(context.TODO(), key, source1, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []string{source1, source2}, sources)

		message, err := eventStore.LookUp(context.TODO(), key, source1)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source1, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []*event.Message{
			event.NewMessage(key, source1, content),
			event.NewMessage(key, source2, content)},
			messageArray)
	})

	t.Run("without folder prefix", func(t *testing.T) {
		bucket := "without-prefix"
		setup(t, bucket)
		config := gcs_event_store.Config(bucket).WithFolder(folderName)
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source1, content)
		assert.NoError(t, err)
		err = eventStore.Persist(context.TODO(), key, source2, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []string{source1, source2}, sources)

		message, err := eventStore.LookUp(context.TODO(), key, source1)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source1, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, []*event.Message{
			event.NewMessage(key, source1, content),
			event.NewMessage(key, source2, content)},
			messageArray)
	})

	t.Run("with compression", func(t *testing.T) {
		bucket := "with-compression"
		setup(t, bucket)
		config := gcs_event_store.Config(bucket).WithCompressor(compression.Gzip(gzip.BestSpeed))
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source1, content)
		assert.NoError(t, err)
		err = eventStore.Persist(context.TODO(), key, source2, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{source1, source2}, sources)

		message, err := eventStore.LookUp(context.TODO(), key, source1)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source1, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []*event.Message{
			event.NewMessage(key, source1, content),
			event.NewMessage(key, source2, content)},
			messageArray)
	})

	t.Run("gcs error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s", key, source1)): // read operation
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				w.Write([]byte(`{}`))
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		millisecond := time.Millisecond
		config := gcs_event_store.Config("bucket").WithTimeout(gcs_event_store.Timeout{Default: &millisecond})
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source1, content)
		assert.Error(t, err)

		_, err = eventStore.ListSourcesByKey(context.TODO(), key)
		assert.Error(t, err)

		_, err = eventStore.LookUp(context.TODO(), key, source1)
		assert.Error(t, err)

		_, err = eventStore.LookUpByKey(context.TODO(), key)
		assert.Error(t, err)
	})
}

func setup(t *testing.T, bucket string) {
	client, err := gcs.NewClient(context.TODO(), internaloption.SkipDialSettingsValidation())
	assert.NoError(t, err)

	_ = client.Bucket(bucket).Create(context.TODO(), "test", nil)
}
