# Event Driver

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lukecold/event-driver)
[![Go Project Version](https://badge.fury.io/go/github.com%2Flukecold%2Fevent-driver.svg)](https://badge.fury.io/go/github.com%2Flukecold%2Fevent-driver)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=lukecold_event-driver&metric=alert_status)](https://sonarcloud.io/dashboard?id=lukecold_event-driver)
![Sonar Coverage](https://img.shields.io/sonar/coverage/lukecold_event-driver?server=https%3A%2F%2Fsonarcloud.io)

Event Driver is a lightweight and flexible event-driven programming framework for managing and handling events in your applications. It provides a simple and intuitive API to facilitate communication between different components or modules in your software.

# Table of Contents
1. [Features](#Features)
2. [Tutorial](#Tutorial)
   1. [Steps Break Down](#Steps-Break-Down)
   2. [Example Code](#Example-Code)
3. [Extensions](#Extensions)
   1. [Cloud Events](#Cloud-Events)
   2. [Google Cloud](#Google-Cloud)

## Features

- **Event-driven Architecture**: Easily implement an event-driven architecture in your application.
- **Custom Handlers**: Define and dispatch custom handlers tailored to your application's needs.
- **Pipeline Structure**: Simply put your handlers in order and expect it to work like a pipeline.
- **Asynchronous Support**: Handle events asynchronously for improved performance and responsiveness.
- **Lightweight and Easy to Use**: Minimal dependencies and quick integration & usage.

## Tutorial

This tutorial is in the format of a case study of a real-world service.

There's an event-driven service that processes orders when `event1` happens,
but the processing logic only uses data from `event2` and `event3`.
This service needs to
1. be fault-tolerant (shouldn't crash at errors, and should be able to recover fast when crash does happen)
2. be idempotent (not processing duplicate orders)

### Steps Break Down
1. Let's start with joining the events. Create a joiner that makes a join when all required events are present.
   ```golang
   myJoiner := joiner.New(joiner.MatchAll("event1", "event2", "event3"), myEventStore)
   ```
2. The event joiner that we just created requires an event store for lookups.
   Here we pick GCS rather than the in-memory store to avoid losing data at restart.
   ```golang
   gcsTimeout := time.Second * 5 // tune the timeout that fits your content size
   gcsConfig := gcs_event_store.Config("my-bucket").
       WithTimeout(gcs_event_store.Timeout{Default: &gcsTimeout}) // can also specify finer timeout for each operation
   // Create a GCS event store without authentication just for showcase.
   myEventStore, err := gcs_event_store.New(ctx, gcsConfig, option.WithoutAuthentication())
   ```
3. Create a cache for idempotency.
   The cache stores the events under the same GCS bucket as joiner (beware of source name conflict between them).

   The cache achieves idempotency by skipping the process (i.e. not passing result to the next handler) on key conflict.
   ```golang
   idempotencyHandler := cache.New(myEventStore, cache.SkipOnConflict())
   ```
4. It would also be great if we can just ignore the content from `event1` to save
   both GCS storage and the network bandwidth & memory of service pod.
   ```golang
   eraseContentOfEvent1 := transformer.EraseContentFromSources("event1")
   eraseContentOfEvent1Handler := transformer.New(eraseContentOfEvent1)
   ```
5. Now we need to build a customized handler to handle the business logic.
   ```golang
   package business_handler

   import (
       "github.com/lukecold/event-driver/handlers"
   )

   type MyInput struct {
       event2 Event2Type `json:"event2"`
       event3 Event3Type `json:"event3"`
   }

   type myHandler struct {
       config Config
   }

   func New(config Config) handlers.Handler {
       return &myHandler{
           config: config,
       }
   }

   func (c *myHandler) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
       var myInput MyInput
       if err := json.UnMarshal(in.GetContent(), &myInput); err != nil {
           return err
       }
       if err = handlesBusinessLogic(myInput, c.config); err != nil {
           return err
       }

       return next.Call(ctx, in) // can omit if this is already the last handler
   }
   ```
6. Build a pipeline with the handlers
   ```golang
   myPipeline := pipeline.New().
       WithNextHandler(eraseContentOfEvent1Handler). // remove the content of event1 before joining
       WithNextHandler(myJoiner).                    // join all events of the same key into a single message
       WithNextHandler(idempotencyHandler)           // check for idempotency after a joint message is formed
   ```
7. Start serving traffic!
   ```golang
   // Showcasing converting the pipeline into KNative cloud events handler
   // One can also use sarama/confluentic kafka, Google Cloud function, etc.
   handleKNativeEvent := convert.ToKNativeEventHandler(
       convert.CloudEventToInput,
       myPipeline,
       convert.OutputToCloudResult)
   cloudEventClient.StartReceiver(ctx, handleKNativeEvent) // assume cloudEventClient is already created
   ```

### Example Code
Now assemble the step break-downs above to an example code

```golang
package main

import (
   "context"
   "log"
   "time"

   "github.com/lukecold/event-driver/event"
   "github.com/lukecold/event-driver/extensions/cloudevents/convert"
   "github.com/lukecold/event-driver/extensions/google-cloud/storage/gcs_event_store"
   "github.com/lukecold/event-driver/handlers/cache"
   "github.com/lukecold/event-driver/handlers/joiner"
   "github.com/lukecold/event-driver/handlers/transformer"
   "github.com/lukecold/event-driver/pipeline"
   "github.com/lukecold/event-driver/storage"
   "google.golang.org/api/option"
)

func main() {
    ctx := context.Background()
    eraseContentOfEvent1 := transformer.EraseContentFromSources("event1")
    eraseContentOfEvent1Handler := transformer.New(eraseContentOfEvent1)
    myEventStore, err := createEventStore(ctx)
    if err != nil {
        log.Panic("failed to create GCS event store", err)
    }
    myJoiner := joiner.New(joiner.MatchAll("event1", "event2", "event3"), myEventStore)
    idempotencyHandler := cache.New(myEventStore, cache.SkipOnConflict())
    businessHandler := business_handler.New()

    myPipeline := pipeline.New().
        WithNextHandler(eraseContentOfEvent1Handler). // remove the content of event1 before joining
        WithNextHandler(myJoiner).                    // join all events of the same key into a single message
        WithNextHandler(idempotencyHandler).          // check for idempotency after a joint message is formed
        WithNextHandler(businessHandler)              // handles the business logic

    // Showcasing converting the pipeline into KNative cloud events handler
    // One can also use sarama/confluentic kafka, Google Cloud function, etc.
    handleKNativeEvent := convert.ToKNativeEventHandler(
        convert.CloudEventToInput,
        myPipeline,
        convert.OutputToCloudResult)
    cloudEventClient.StartReceiver(ctx, handleKNativeEvent) // assume cloudEventClient is already created
}

func createEventStore(ctx context.Context) (storage.EventStore, error) {
    gcsTimeout := time.Second * 5 // tune the timeout that fits your content size
    gcsConfig := gcs_event_store.Config("my-bucket").
        WithTimeout(gcs_event_store.Timeout{Default: &gcsTimeout}) // can also specify finer timeout for each operation

    // Create a GCS event store without authentication just for showcase.
    return gcs_event_store.New(ctx, gcsConfig, option.WithoutAuthentication())
}
```

## Extensions
Event Driver also provides the following libs for integrating with other services/frameworks as extensions.

### Cloud Events

Link: [github.com/lukecold/event-driver/extensions/cloudevents](https://github.com/lukecold/event-driver/tree/main/extensions/cloudevents)

Integrate event driver with [Cloud Events](github.com/cloudevents/sdk-go/v2),
i.e. providing a converter that converts the event driver pipeline into a cloud events handler.
Check the [document](https://github.com/lukecold/event-driver/tree/main/extensions/cloudevents/README.md)
to see what is currently supported and the latest update.

### Google Cloud

Link: [github.com/lukecold/event-driver/extensions/google-cloud](https://github.com/lukecold/event-driver/tree/main/extensions/google-cloud)

Integrate event driver with Google Cloud, including using GCS/BigQuery as event store,
integrating the event driver pipeline with Cloud Functions, etc.
Check the [document](https://github.com/lukecold/event-driver/tree/main/extensions/google-cloud/README.md)
to see what is currently supported and the latest update.
