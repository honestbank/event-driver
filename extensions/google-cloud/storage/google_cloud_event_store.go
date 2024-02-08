package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
)

// GoogleCloudEventStore persists the contents in GCS, which requires consistent connections to Google Cloud.
type GoogleCloudEventStore struct {
	bucketName string
	client     *gcs.Client
}

type GCSConfig struct {
	AccessToken string
	Keyfile     string
	Bucket      string
}

func NewGoogleCloudEventStore(ctx context.Context, cfg GCSConfig) (storage.EventStore, error) {
	token := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: cfg.AccessToken,
	})

	client, err := gcs.NewClient(ctx, option.WithTokenSource(token), internaloption.SkipDialSettingsValidation())
	if err != nil {
		return nil, err
	}

	return &GoogleCloudEventStore{
		bucketName: cfg.Bucket,
		client:     client,
	}, nil
}

func (g *GoogleCloudEventStore) Persist(ctx context.Context, key, source, content string) error {
	reqCtx, cancelFunc := context.WithTimeout(ctx, time.Second*50)
	defer cancelFunc()
	bucket := g.client.Bucket(g.bucketName)

	filename := fmt.Sprintf("%s/%s", key, source)
	wc := bucket.Object(filename).NewWriter(reqCtx)
	if _, err := io.Copy(wc, io.NopCloser(bytes.NewBufferString(content))); err != nil {
		return err
	}

	return wc.Close()
}

func (g *GoogleCloudEventStore) LookUp(ctx context.Context, key, source string) (*event.Message, error) {
	reqCtx, cancelFunc := context.WithTimeout(ctx, time.Second*50)
	defer cancelFunc()
	bucket := g.client.Bucket(g.bucketName)

	filename := fmt.Sprintf("%s/%s", key, source)
	content, err := g.readFile(reqCtx, bucket, filename)
	if err != nil {
		return nil, err
	}

	return event.NewMessage(key, source, string(content)), nil
}

func (g *GoogleCloudEventStore) LookUpByKey(ctx context.Context, key string) ([]*event.Message, error) {
	reqCtx, cancelFunc := context.WithTimeout(ctx, time.Second*50)
	defer cancelFunc()
	bucket := g.client.Bucket(g.bucketName)

	messages := make([]*event.Message, 0)
	objectIterator := bucket.Objects(reqCtx, &gcs.Query{Prefix: key + "/"})
	for {
		object, err := objectIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return messages, err
		}
		content, err := g.readFile(reqCtx, bucket, object.Name)
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

func (g *GoogleCloudEventStore) readFile(ctx context.Context, bucket *gcs.BucketHandle, filename string) ([]byte, error) {
	reader, err := bucket.Object(filename).NewReader(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, nil
		}

		return nil, err
	}
	content := make([]byte, 0)
	_, err = reader.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func parseFileName(filename string) (key, source string, err error) {
	split := strings.Split(filename, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("filename %s isn't of format 'key/source'", filename)
	}

	return split[0], split[1], nil
}
