package gcs_event_store_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/extensions/google-cloud/storage/gcs_event_store"
)

const (
	key     = "test-key"
	source  = "test-source"
	content = "test-content"
)

func TestGCSEventStore(t *testing.T) {
	t.Run("with folder prefix", func(t *testing.T) {
		folderName := "folder-name"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Println(r.URL.Path)
			switch {
			case strings.HasPrefix(r.URL.Path, "/upload/"): // write operation
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				// assert that bucket, key, source, and content match
				assert.Contains(t, string(body), fmt.Sprintf(`{"bucket":"bucket","name":"%s/%s/%s"}`, folderName, key, source))
				assert.Contains(t, string(body), content)
				w.Write([]byte("{}"))
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s/%s", folderName, key, source)): // read operation
				w.Write([]byte(content))
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				w.Write([]byte(`{
				"kind":"storage#objects",
				"items":[{"kind":"storage#object","name":"folderName/key/source1"},{"kind":"storage#object","name":"folderName/key/source2"}]
			}`))
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		config := gcs_event_store.Config("bucket").WithFolder(folderName)
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.NoError(t, err)

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
				assert.Contains(t, string(body), fmt.Sprintf(`{"bucket":"bucket","name":"%s/%s"}`, key, source))
				assert.Contains(t, string(body), content)
				w.Write([]byte("{}"))
			case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/bucket/%s/%s", key, source)): // read operation
				w.Write([]byte(content))
			case strings.HasPrefix(r.URL.Path, "/storage/v1/b/bucket/o"): // list operation
				w.Write([]byte(`{
				"kind":"storage#objects",
				"items":[{"kind":"storage#object","name":"key/source1"},{"kind":"storage#object","name":"key/source2"}]
			}`))
			}
		}))
		defer server.Close()

		os.Setenv("STORAGE_EMULATOR_HOST", server.URL)

		config := gcs_event_store.Config("bucket")
		eventStore, err := gcs_event_store.New(context.TODO(), config, option.WithoutAuthentication())
		assert.NoError(t, err)

		err = eventStore.Persist(context.TODO(), key, source, content)
		assert.NoError(t, err)

		message, err := eventStore.LookUp(context.TODO(), key, source)
		assert.NoError(t, err)
		assert.Equal(t, event.NewMessage(key, source, content), message)

		messageArray, err := eventStore.LookUpByKey(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(messageArray))
	})
}
