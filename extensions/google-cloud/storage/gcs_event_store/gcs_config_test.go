package gcs_event_store_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/event-driver/extensions/google-cloud/storage/gcs_event_store"
)

func TestGCSConfig(t *testing.T) {
	defaultTimeout := time.Minute
	listTimeout := time.Second * 10
	readTimeout := time.Second * 20
	writeTimeout := time.Second * 30

	operationToTimeout := map[gcs_event_store.Operation]time.Duration{
		gcs_event_store.ListContents: listTimeout,
		gcs_event_store.ReadContent:  readTimeout,
		gcs_event_store.WriteContent: writeTimeout,
	}
	gcsConfig := gcs_event_store.Config("bucket").
		WithTimeout(gcs_event_store.Timeout{
			Default:   &defaultTimeout,
			Operation: operationToTimeout,
		})

	// use timeout of specific operation
	for operation, expectedTimeout := range operationToTimeout {
		t.Run(fmt.Sprintf("operation %s", operation), func(t *testing.T) {
			timeBefore := time.Now()
			ctx, _ := gcsConfig.NewContextWithTimeout(context.Background(), operation)
			timeAfter := time.Now()
			deadline, isDeadlineConfigured := ctx.Deadline()
			assert.True(t, isDeadlineConfigured)
			assert.True(t, timeBefore.Add(expectedTimeout).Before(deadline))
			assert.True(t, timeAfter.Add(expectedTimeout).After(deadline))
		})
	}

	// use default timeout if operation isn't configured
	operationNotConfigured := gcs_event_store.Operation("NotConfigured")
	t.Run("operation not configured", func(t *testing.T) {
		timeBefore := time.Now()
		ctx, _ := gcsConfig.NewContextWithTimeout(context.Background(), operationNotConfigured)
		timeAfter := time.Now()
		deadline, isDeadlineConfigured := ctx.Deadline()
		assert.True(t, isDeadlineConfigured)
		assert.True(t, timeBefore.Add(defaultTimeout).Before(deadline))
		assert.True(t, timeAfter.Add(defaultTimeout).After(deadline))
	})

	// when timeout isn't configured, use a universal timeout of 30s
	cfgNoTimeout := gcs_event_store.Config("bucket")
	t.Run("no timeout", func(t *testing.T) {
		timeBefore := time.Now()
		ctx, _ := cfgNoTimeout.NewContextWithTimeout(context.Background(), operationNotConfigured)
		timeAfter := time.Now()
		deadline, isDeadlineConfigured := ctx.Deadline()
		assert.True(t, isDeadlineConfigured)
		assert.True(t, timeBefore.Add(30*time.Second).Before(deadline))
		assert.True(t, timeAfter.Add(30*time.Second).After(deadline))
	})
}
