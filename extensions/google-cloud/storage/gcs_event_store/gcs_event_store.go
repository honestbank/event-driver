package gcs_event_store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/storage"
)

// GCSEventStore persists the contents in GCS, which requires consistent connections to Google Cloud.
type GCSEventStore struct {
	cfg    *GCSConfig
	client *gcs.Client
}

func New(ctx context.Context, cfg *GCSConfig, options ...option.ClientOption) (storage.EventStore, error) {
	if cfg == nil {
		return nil, errors.New("gcs config cannot be null")
	}
	client, err := gcs.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}

	return &GCSEventStore{
		cfg:    cfg,
		client: client,
	}, nil
}

// Persist uploads the message as a file on the path `key/source`.
func (g *GCSEventStore) Persist(ctx context.Context, key, source, content string) error {
	bucket := g.client.Bucket(g.cfg.Bucket)

	filename := fmt.Sprintf("%s/%s", key, source)

	writeRequestCtx, _ := g.cfg.NewContextWithTimeout(ctx, WriteContent)
	writer := bucket.Object(filename).NewWriter(writeRequestCtx)
	if _, err := writer.Write([]byte(content)); err != nil {
		return err
	}

	return writer.Close()
}

// LookUp returns a single message by looking up the path `key/source`.
func (g *GCSEventStore) LookUp(ctx context.Context, key, source string) (*event.Message, error) {
	bucket := g.client.Bucket(g.cfg.Bucket)

	filename := fmt.Sprintf("%s/%s", key, source)
	readRequestCtx, _ := g.cfg.NewContextWithTimeout(ctx, ReadContent)
	content, err := g.readFile(readRequestCtx, bucket, filename)
	if err != nil {
		return nil, err
	}

	return event.NewMessage(key, source, string(content)), nil
}

// LookUpByKey returns a list of messages by looking up the prefix `key/`.
func (g *GCSEventStore) LookUpByKey(ctx context.Context, key string) ([]*event.Message, error) {
	bucket := g.client.Bucket(g.cfg.Bucket)

	messages := make([]*event.Message, 0)
	listRequestCtx, _ := g.cfg.NewContextWithTimeout(ctx, ListContents)
	objectIterator := bucket.Objects(listRequestCtx, &gcs.Query{Prefix: key + "/"})
	for {
		object, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return messages, err
		}
		readRequestCtx, _ := g.cfg.NewContextWithTimeout(ctx, ReadContent)
		content, err := g.readFile(readRequestCtx, bucket, object.Name)
		if err != nil {
			return messages, err
		}
		key, source, err := parseFileName(object.Name)
		if err != nil {
			return messages, err
		}
		messages = append(messages, event.NewMessage(key, source, string(content)))
	}

	return messages, nil
}

func (g *GCSEventStore) readFile(ctx context.Context, bucket *gcs.BucketHandle, filename string) ([]byte, error) {
	reader, err := bucket.Object(filename).NewReader(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, nil
		}

		return nil, err
	}

	return io.ReadAll(reader)
}

func parseFileName(filename string) (key, source string, err error) {
	split := strings.Split(filename, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("filename %s isn't of format 'key/source'", filename)
	}

	return split[0], split[1], nil
}
