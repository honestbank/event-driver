# Event Driver - cloud events extension

## Construction Checklist
- [x] Support cloud events vended by KNative kafka-source
- [ ] Create a feature-request or pull-request if you need something more

## Usage

The following showcases an example of building a message processing pipeline, and convert it into a cloud events handler.

```golang
package main

import (
    "context"
    "log"

    "github.com/honestbank/event-driver/event"
    "github.com/honestbank/event-driver/extensions/cloudevents/convert"
    "github.com/honestbank/event-driver/handlers/cache"
    "github.com/honestbank/event-driver/handlers/joiner"
    "github.com/honestbank/event-driver/handlers/transformer"
    "github.com/honestbank/event-driver/pipeline"
    "github.com/honestbank/event-driver/storage"
)

func main() {
    ctx := context.Background()
    renameSources, err := transformer.RenameSources(map[string][]string{"source1": {"alias1", "alias2"}})
    if err != nil {
        log.Panic("failed to create 'RenameSources' transformer", err)
    }
    myPipeline := pipeline.New().
        WithNextHandler(renameSources).
        WithNextHandler(joiner.New(joiner.MatchAll("source1", "source2"), storage.NewInMemoryStore())).
        WithNextHandler(cache.New(storage.NewInMemoryStore(), cache.SkipOnConflict()))

    // if convert to cloud events handler
    handleKNativeEvent := convert.ToKNativeEventHandler(
        convert.CloudEventToInput,
        myPipeline,
        convert.OutputToCloudResult)
    cloudEventClient.StartReceiver(ctx, handleKNativeEvent)
}
```
