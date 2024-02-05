package convert

import (
	"context"
	"errors"
	"net/http"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/pipeline"
)

func ToKNativeEventHandler(
	convertInput func(cloudevents.Event) (*event.Message, error),
	pipeline pipeline.Pipeline,
	convertOutput func(error) (*cloudevents.Event, cloudevents.Result)) func(
	ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	return func(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
		input, err := convertInput(event)
		if err != nil {
			return nil, cloudevents.NewHTTPResult(http.StatusBadRequest, "%s", err)
		}
		output := pipeline.Run(ctx, input)

		return convertOutput(output)
	}
}

func GetKey(event cloudevents.Event) (*string, error) {
	rawKey := event.Extensions()["key"]
	if rawKey == nil {
		return nil, errors.New("cannot find key from event extension")
	}
	key, isString := rawKey.(string)
	if !isString {
		return nil, errors.New("key isn't of string type")
	}

	return &key, nil
}

func GetTopic(event cloudevents.Event) (*string, error) {
	source := event.Source()
	split := strings.Split(source, "#")
	if len(split) == 0 {
		return nil, errors.New("cannot find topic from event source")
	}

	return &split[len(split)-1], nil
}
