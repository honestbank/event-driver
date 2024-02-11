package convert

import (
	"context"
	"errors"
	"net/http"
	"strings"

	cloudEvents "github.com/cloudevents/sdk-go/v2"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/pipeline"
)

// InputConverter converts the KNativeEventHandler input to pipeline.Pipeline input.
type InputConverter func(context.Context, *cloudEvents.Event) (*event.Message, error)

// OutputConverter converts the pipeline.Pipeline output to KNativeEventHandler output.
type OutputConverter func(context.Context, error) (*cloudEvents.Event, cloudEvents.Result)

// KNativeEventHandler represents the most complicated function to
// [start an event receiver](https://github.com/cloudevents/sdk-go/blob/main/v2/client/client.go#L33-L50),
// in order to provide compatibility for all use cases.
type KNativeEventHandler func(context.Context, cloudEvents.Event) (*cloudEvents.Event, cloudEvents.Result)

// ToKNativeEventHandler converts the pipeline to a KNativeEventHandler,
// with an InputConverter to convert the KNativeEventHandler input to pipeline.Pipeline input,
// and an OutputConverter to convert the pipeline.Pipeline output to KNativeEventHandler output.
func ToKNativeEventHandler(
	convertInput InputConverter,
	pipeline pipeline.Pipeline,
	convertOutput OutputConverter) KNativeEventHandler {
	return func(ctx context.Context, cloudEvent cloudEvents.Event) (*cloudEvents.Event, cloudEvents.Result) {
		input, err := convertInput(ctx, &cloudEvent)
		if err != nil {
			return nil, cloudEvents.NewHTTPResult(http.StatusBadRequest, "%s", err)
		}
		output := pipeline.Process(ctx, input)

		return convertOutput(ctx, output)
	}
}

// CloudEventToInput is a built-in InputConverter that could serve the basic functionalities of
// converting a cloud event to a pipeline input message.
func CloudEventToInput(_ context.Context, cloudEvent *cloudEvents.Event) (*event.Message, error) {
	key, err := GetKey(cloudEvent)
	if err != nil {
		return nil, err
	}
	source, err := GetTopic(cloudEvent)
	if err != nil {
		return nil, err
	}

	return event.NewMessage(*key, *source, string(cloudEvent.Data())), nil
}

// OutputToCloudResult is a built-in OutputConverter that could serve the basic functionalities of
// converting an error (the pipeline output) to the cloud event response.
func OutputToCloudResult(_ context.Context, err error) (*cloudEvents.Event, cloudEvents.Result) {
	if err != nil {
		return nil, cloudEvents.NewHTTPResult(http.StatusInternalServerError, "%s", err)
	}

	return nil, cloudEvents.NewHTTPResult(http.StatusOK, "OK")
}

// GetKey tries to find the event key of string format under key `key`, in the extensions map.
func GetKey(cloudEvent *cloudEvents.Event) (*string, error) {
	rawKey := cloudEvent.Extensions()["key"]
	if rawKey == nil {
		return nil, errors.New("cannot find key from event extension")
	}
	key, isString := rawKey.(string)
	if !isString {
		return nil, errors.New("key isn't of string type")
	}

	return &key, nil
}

// GetTopic tries to find the topic that's stored as event source.
func GetTopic(cloudEvent *cloudEvents.Event) (*string, error) {
	source := cloudEvent.Source()
	if len(source) == 0 {
		return nil, errors.New("event source is empty")
	}
	split := strings.Split(source, "#")
	if len(split) == 0 {
		return nil, errors.New("cannot find topic from event source")
	}

	return &split[len(split)-1], nil
}
