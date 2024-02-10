package convert_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	cloudEvents "github.com/cloudevents/sdk-go/v2/event"
	cloudEventTypes "github.com/cloudevents/sdk-go/v2/types"
	"github.com/heetch/avro"
	"github.com/stretchr/testify/assert"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/extensions/cloudevents/convert"
	"github.com/lukecold/event-driver/pipeline"
)

const (
	key   = "test-key"
	topic = "test-topic"
)

type eventContent struct {
	StringVal string
	IntVal    int
	FloatPtr  *float64
}

func TestGetKey(t *testing.T) {
	t.Run("get key of string type", func(t *testing.T) {
		cloudEvent := makeEvent(key, topic, nil)
		parsedKey, err := convert.GetKey(cloudEvent)
		assert.NoError(t, err)
		assert.Equal(t, key, *parsedKey)
	})

	t.Run("fail if key cannot be casted to string", func(t *testing.T) {
		cloudEvent := makeEvent(1, topic, nil)
		_, err := convert.GetKey(cloudEvent)
		assert.Error(t, err)
	})

	t.Run("failed if key isn't present in extension map", func(t *testing.T) {
		cloudEvent := &cloudEvents.Event{
			Context: cloudEvents.EventContextV03{
				Source:          cloudEventTypes.URIRef{URL: url.URL{Fragment: topic}},
				DataContentType: cloudEvents.StringOfApplicationJSON(),
				Extensions:      map[string]interface{}{},
			}.AsV03(),
			DataEncoded: nil,
		}
		_, err := convert.GetKey(cloudEvent)
		assert.Error(t, err)
	})
}

func TestGetTopic(t *testing.T) {
	t.Run("get topic from source", func(t *testing.T) {
		cloudEvent := makeEvent(key, topic, nil)
		parsedTopic, err := convert.GetTopic(cloudEvent)
		assert.NoError(t, err)
		assert.Equal(t, topic, *parsedTopic)
	})

	t.Run("fail if source is empty", func(t *testing.T) {
		cloudEvent := makeEvent(key, "", nil)
		_, err := convert.GetTopic(cloudEvent)
		assert.Error(t, err)
	})
}

func TestCloudEventToInput(t *testing.T) {
	floatVal := 1.1
	content, err := json.Marshal(eventContent{
		StringVal: "test-string-value",
		IntVal:    123,
		FloatPtr:  &floatVal,
	})
	assert.NoError(t, err)
	cloudEvent := makeEvent(key, topic, content)
	input, err := convert.CloudEventToInput(context.TODO(), cloudEvent)
	assert.NoError(t, err)
	expectedInput := event.NewMessage(key, topic, string(content))
	assert.Equal(t, expectedInput, input)
}

func TestCloudEventAVROToInput(t *testing.T) {
	avroData := []byte(`{"id":"123","source":"example.com","specversion":"1.0","type":"eventType","datacontenttype":"application/json","data":"{\"key\":\"value\"}"}`)
	stringOfTypeAVRO := "avro/binary"

	cloudEvent := &v2.Event{
		Context: &v2.EventContextV03{
			Type:            "com.example.test",
			Source:          cloudEventTypes.URIRef{URL: url.URL{Fragment: topic}},
			ID:              "id123",
			DataContentType: &stringOfTypeAVRO,
		},
		DataEncoded: avroData,
	}

	// Mock schemaType for AVRO unmarshaling
	schemaType := &avro.Type{} // You would provide an actual schemaType here

	// Call the function
	message, err := convert.CloudEventAVROToInput(context.Background(), cloudEvent, schemaType)

	// Check if the result is as expected
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "id123", message.GetKey())
	assert.Equal(t, "http://example.com/source", message.GetSource())
	assert.Equal(t, "{\"key\":\"value\"}", message.GetContent()) // This assumes you are storing the payload as a string
}

// A plain happy case of convert.ToKNativeEventHandler.
func TestToKNativeEventHandler(t *testing.T) {
	emptyPipeline := pipeline.New()
	handleKNativeEvent := convert.ToKNativeEventHandler(
		convert.CloudEventToInput,
		emptyPipeline,
		convert.OutputToCloudResult)
	responseEvent, result := handleKNativeEvent(context.TODO(), *makeEvent(key, topic, nil))
	assert.Nil(t, responseEvent)
	assert.Equal(t, v2.NewHTTPResult(http.StatusOK, "OK"), result)
}

func makeEvent(key interface{}, topic string, content []byte) *cloudEvents.Event {
	return &cloudEvents.Event{
		Context: cloudEvents.EventContextV03{
			Source:          cloudEventTypes.URIRef{URL: url.URL{Fragment: topic}},
			DataContentType: cloudEvents.StringOfApplicationJSON(),
			Extensions: map[string]interface{}{
				"key": key,
			},
		}.AsV03(),
		DataEncoded: content,
	}
}
