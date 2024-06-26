package gcs_event_store

import (
	"context"
	"time"

	"github.com/honestbank/event-driver/utils/compression"
)

type Operation string

const (
	ListContents Operation = "ListContents" // operation that lists all content associated with a given key
	ReadContent  Operation = "ReadContent"  // operation that reads content associated with a key-source pair
	WriteContent Operation = "WriteContent" //operation to writes content associated with a key-source pair
)

type GCSConfig struct {
	Bucket     string
	Compressor compression.Compressor
	Folder     *string
	ReadPolicy ReadPolicy
	Timeout    Timeout
}

type Timeout struct {
	Default   *time.Duration              // the default timeout for all operations
	Operation map[Operation]time.Duration // the timeout of each operation - this overrides the default timeout
}

// Config creates a default configuration, that
// - doesn't do compression/decompression when write & read to GCS
// - takes the earliest created object if there are multiple under the same key/source/ path
// - enforces universal 30s timeout in GCS requests
func Config(bucket string) *GCSConfig {
	halfMinute := 30 * time.Second

	return &GCSConfig{
		Bucket:     bucket,
		Compressor: compression.Noop(),
		ReadPolicy: TakeFirstCreated(),
		Timeout: Timeout{
			Default:   &halfMinute,
			Operation: make(map[Operation]time.Duration),
		},
	}
}

func (c *GCSConfig) WithCompressor(compressor compression.Compressor) *GCSConfig {
	c.Compressor = compressor

	return c
}

func (c *GCSConfig) WithFolder(folder string) *GCSConfig {
	c.Folder = &folder

	return c
}

func (c *GCSConfig) WithReadPolicy(readPolicy ReadPolicy) *GCSConfig {
	c.ReadPolicy = readPolicy

	return c
}

func (c *GCSConfig) WithTimeout(timeout Timeout) *GCSConfig {
	c.Timeout = timeout

	return c
}

// NewContextWithTimeout generates a new context.Context with timeout from the parent context.
// If neither operation-specific timeout nor default timeout is found, return parent context without any operation.
func (c *GCSConfig) NewContextWithTimeout(
	parent context.Context,
	operation Operation) (context.Context, context.CancelFunc) {
	if listContentsTimeout, isConfigured := c.Timeout.Operation[operation]; isConfigured {
		return context.WithTimeout(parent, listContentsTimeout)
	}
	if c.Timeout.Default != nil {
		return context.WithTimeout(parent, *c.Timeout.Default)
	}

	return parent, noop
}

func noop() {
}
