package gcs_event_store_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/extensions/google-cloud/storage/gcs_event_store"
	"github.com/honestbank/event-driver/utils/compression"
)

const (
	key     = "key"
	source  = "source"
	content = "content"
)

func TestGCSEventStore(t *testing.T) {
	t.Run("with folder prefix", func(t *testing.T) {
		folderName := "folder-name"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				// assert that bucket, key, source, and content match
				assert.Contains(t, string(body), fmt.Sprintf(`{"bucket":"bucket","name":"%s/%s/%s`, folderName, key, source))
				assert.Contains(t, string(body), content)
				w.Write([]byte("{}"))
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s/%s", folderName, key, source)): // read operation
				w.Write([]byte(content))
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				if strings.Contains(r.RequestURI, "delimiter=%2F") {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"folder-name/key/source1"},{"kind":"storage#object","name":"folder-name/key/source2"}]
					}`))
				} else {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"folder-name/key/source1/sha1"},{"kind":"storage#object","name":"folder-name/key/source1/sha2"}]
					}`))
				}
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		config := gcs_event_store.Config("bucket").WithFolder(folderName)
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(sources))

		message, err := eventStore.LookUp(context.TODO(), key, source)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(messageArray))
	})

	t.Run("without folder prefix", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Println(r.URL.Path)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				// assert that bucket, key, source, and content match
				assert.Contains(t, string(body), fmt.Sprintf(`{"bucket":"bucket","name":"%s/%s`, key, source))
				assert.Contains(t, string(body), content)
				w.Write([]byte("{}"))
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s", key, source)): // read operation
				w.Write([]byte(content))
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				if strings.Contains(r.RequestURI, "delimiter=%2F") {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"key/source1"},{"kind":"storage#object","name":"key/source2"}]
					}`))
				} else {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"key/source1/sha1"},{"kind":"storage#object","name":"key/source1/sha2"}]
					}`))
				}
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		config := gcs_event_store.Config("bucket")
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(sources))

		message, err := eventStore.LookUp(context.TODO(), key, source)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(messageArray))
	})

	t.Run("with compression", func(t *testing.T) {
		compressor := compression.Gzip(gzip.BestSpeed)
		compressedContent, err := compressor.Compress([]byte(content))
		assert.NoError(t, err)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				// assert that bucket, key, source, and content match
				assert.Contains(t, string(body), fmt.Sprintf(`{"bucket":"bucket","name":"%s/%s`, key, source))
				assert.Contains(t, string(body), string(compressedContent))
				w.Write([]byte("{}"))
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s", key, source)): // read operation
				w.Write(compressedContent)
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				if strings.Contains(r.RequestURI, "delimiter=%2F") {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"key/source1"},{"kind":"storage#object","name":"key/source2"}]
					}`))
				} else {
					w.Write([]byte(`{
						"kind":"storage#objects",
						"items":[{"kind":"storage#object","name":"key/source1/sha1"},{"kind":"storage#object","name":"key/source1/sha2"}]
					}`))
				}
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		config := gcs_event_store.Config("bucket").WithCompressor(compression.Gzip(gzip.BestSpeed))
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.NoError(t, err)

		sources, err := eventStore.ListSourcesByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(sources))

		message, err := eventStore.LookUp(context.TODO(), key, source)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(messageArray))
	})

	t.Run("gcs error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s", key, source)): // read operation
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

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.Error(t, err)

		_, err = eventStore.ListSourcesByKey(context.TODO(), key)
		assert.Error(t, err)

		_, err = eventStore.LookUp(context.TODO(), key, source)
		assert.Error(t, err)

		_, err = eventStore.LookUpByKey(context.TODO(), key)
		assert.Error(t, err)
	})
}
