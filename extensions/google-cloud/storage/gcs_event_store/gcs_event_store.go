package gcs_event_store

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/storage"
	"github.com/lukecold/event-driver/utils/compression"
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

func (g *GCSEventStore) ListSourcesByKey(ctx context.Context, key string) ([]string, error) {
	bucket := g.client.Bucket(g.cfg.Bucket)

	sources := make([]string, 0)
	listRequestCtx, cancel := g.cfg.NewContextWithTimeout(ctx, ListContents)
	defer cancel()

	objectIterator := bucket.Objects(listRequestCtx, &gcs.Query{
		Prefix:    composePath(g.cfg.Folder, key) + "/",
		Delimiter: "/",
	})
	for {
		object, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return sources, err
		}
		_, source, err := parsePath(object.Name)
		if err != nil {
			return sources, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

// LookUp returns a single message by looking up the path `folder/key/source`.
func (g *GCSEventStore) LookUp(ctx context.Context, key, source string) (*event.Message, error) {
	bucket := g.client.Bucket(g.cfg.Bucket)

	path := composePath(g.cfg.Folder, key, source)
	readRequestCtx, cancel := g.cfg.NewContextWithTimeout(ctx, ReadContent)
	defer cancel()

	content, err := readFile(readRequestCtx, g.cfg.Compressor, bucket, path, g.cfg.ReadPolicy)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil
	}

	return event.NewMessage(key, source, string(content)), nil
}

// LookUpByKey returns a list of messages by looking up the prefix `folder/key/`.
func (g *GCSEventStore) LookUpByKey(ctx context.Context, key string) ([]*event.Message, error) {
	bucket := g.client.Bucket(g.cfg.Bucket)

	messages := make([]*event.Message, 0)
	listRequestCtx, cancel := g.cfg.NewContextWithTimeout(ctx, ListContents)
	defer cancel()

	sourceIterator := bucket.Objects(listRequestCtx, &gcs.Query{
		Prefix:    composePath(g.cfg.Folder, key) + "/",
		Delimiter: "/",
	})
	for {
		path, err := sourceIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return messages, err
		}
		readRequestCtx, _ := g.cfg.NewContextWithTimeout(ctx, ReadContent)
		content, err := readFile(readRequestCtx, g.cfg.Compressor, bucket, path.Name, g.cfg.ReadPolicy)
		if err != nil {
			return messages, err
		}
		key, source, err := parsePath(path.Name)
		if err != nil {
			return messages, err
		}
		messages = append(messages, event.NewMessage(key, source, string(content)))
	}

	return messages, nil
}

// Persist uploads the message as a file on the path `folder/key/source`.
func (g *GCSEventStore) Persist(ctx context.Context, key, source, content string) error {
	bucket := g.client.Bucket(g.cfg.Bucket)
	path := composePath(g.cfg.Folder, key, source)
	writeRequestCtx, cancel := g.cfg.NewContextWithTimeout(ctx, WriteContent)
	defer cancel()

	return writeFile(writeRequestCtx, g.cfg.Compressor, bucket, path, []byte(content))
}

func composePath(folder *string, keys ...string) string {
	components := make([]string, 0)

	if folder != nil {
		components = append(components, strings.Trim(*folder, "/"))
	}
	for _, key := range keys {
		components = append(components, strings.Trim(key, "/"))
	}

	return strings.Join(components, "/")
}

func parsePath(filename string) (key, source string, err error) {
	split := strings.Split(filename, "/")
	if len(split) < 2 {
		return "", "", fmt.Errorf("filename %s isn't of format 'folder/key/source'", filename)
	}

	return split[len(split)-2], split[len(split)-1], nil
}

func readFile(
	ctx context.Context,
	compressor compression.Compressor,
	bucket *gcs.BucketHandle,
	path string,
	readPolicy ReadPolicy) ([]byte, error) {
	objectIterator := bucket.Objects(ctx, &gcs.Query{Prefix: path + "/"})
	object, err := readPolicy.Apply(objectIterator)
	if err != nil {
		return nil, err
	}
	reader, err := bucket.Object(object.Name).NewReader(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, nil
		}

		return nil, err
	}
	compressedContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return compressor.Decompress(compressedContent)
}

func writeFile(
	ctx context.Context,
	compressor compression.Compressor,
	bucket *gcs.BucketHandle,
	path string,
	content []byte) error {
	compressedContent, err := compressor.Compress(content)
	if err != nil {
		return err
	}
	sha := sha256.Sum256(compressedContent)
	filename := fmt.Sprintf("%s/%s", path, base64.URLEncoding.EncodeToString(sha[:]))
	writer := bucket.Object(filename).NewWriter(ctx)
	if _, err = writer.Write(compressedContent); err != nil {
		return err
	}

	return writer.Close()
}
