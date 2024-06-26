# Event Driver - Google Cloud extension

## Construction Checklist
- [x] Support GCS event store
- [ ] Support BigQuery event store
- [ ] Integrate event-driver pipeline with Cloud Function
- [ ] Create a feature-request or pull-request if you need something more

## Usage

The following showcases an example of building a message processing pipeline with Google Cloud Storage.

```golang
package main

import (
    "context"
    "log"

    "github.com/honestbank/event-driver/event"
    "github.com/honestbank/event-driver/extensions/cloudevents/convert"
    "github.com/honestbank/event-driver/extensions/google-cloud/storage/gcs_event_store"
    "github.com/honestbank/event-driver/handlers/cache"
    "github.com/honestbank/event-driver/handlers/joiner"
    "github.com/honestbank/event-driver/handlers/transformer"
    "github.com/honestbank/event-driver/pipeline"
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
