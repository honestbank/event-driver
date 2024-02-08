package gcs_event_store

import (
	"context"
	"time"
)

type Operation string

const (
	ListContents Operation = "ListContents" // operation that lists all content associated with a given key
	ReadContent  Operation = "ReadContent"  // operation that reads content associated with a key-source pair
	WriteContent Operation = "WriteContent" //operation to writes content associated with a key-source pair
)

type GCSConfig struct {
	Bucket  string
	Timeout Timeout
}

type Timeout struct {
	Default   *time.Duration              // the default timeout for all operations
	Operation map[Operation]time.Duration // the timeout of each operation - this overrides the default timeout
}

func Config(bucket string) *GCSConfig {
	return &GCSConfig{
		Bucket: bucket,
		Timeout: Timeout{
			Default:   nil,
			Operation: make(map[Operation]time.Duration),
		},
	}
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
