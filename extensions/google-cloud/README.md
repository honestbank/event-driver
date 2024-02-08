# Event Driver - Google Cloud extension

## Usage

The following showcases an example of building a message processing pipeline with Google Cloud Storage.

```golang
package main

import (
	"context"
	"log"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/extensions/google-cloud/storage/gcs_event_store"
	"github.com/lukecold/event-driver/extensions/knative/convert"
	"github.com/lukecold/event-driver/handlers/cache"
	"github.com/lukecold/event-driver/handlers/joiner"
	"github.com/lukecold/event-driver/handlers/transformer"
	"github.com/lukecold/event-driver/pipeline"
)

func main() {
	ctx := context.Background()
	renameSources, err := transformer.RenameSources(map[string][]string{"source1": {"alias1", "alias2"}})
	if err != nil {
		log.Panic("failed to create 'RenameSources' transformer", err)
	}
	gcsEventStore, err := gcs_event_store.New(ctx, gcs_event_store.Config("my-bucket"))
	if err != nil {
		log.Panic("failed to create GCS event store", err)
	}
	myPipeline := pipeline.New().
		WithNextHandler(renameSources).
		WithNextHandler(joiner.New(joiner.MatchAll("source1", "source2"), gcsEventStore)).
		WithNextHandler(cache.New(gcsEventStore, cache.SkipOnConflict()))
	// handle events with myPipeline
}
```
