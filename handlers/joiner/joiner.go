package joiner

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/handlers/options"
	"github.com/lukecold/event-driver/storage"
)

// joiner implements handlers.Handler that joins the events with the same key
// when the sources match the criteria given by Condition.
// The output is of JSON format `{"source1":"content1","source2":"content2",...}`.
type joiner struct {
	condition Condition
	storage   storage.EventStore
	logger    *slog.Logger
}

func New(condition Condition, storage storage.EventStore, opts ...options.Option) *joiner {
	cfg := options.DefaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &joiner{
		condition: condition,
		storage:   storage,
		logger: slog.New(slog.NewJSONHandler(cfg.GetLogWriter(), &slog.HandlerOptions{Level: cfg.GetLogLevel()})).
			With(slog.String("handler", "joiner")),
	}
}

func (j *joiner) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	logger := j.logger.With(slog.String("key", in.GetKey()), slog.String("source", in.GetSource()))
	// persist input message by key & source
	err := j.storage.Persist(ctx, in.GetKey(), in.GetSource(), in.GetContent())
	if err != nil {
		logger.Error("failed to persist message", slog.Any("error", err))

		return err
	}

	// validate sources
	messages, err := j.storage.LookUpByKey(ctx, in.GetKey())
	if err != nil {
		logger.Error("failed to look up by key", slog.Any("error", err))

		return err
	}
	persistedSources := make([]string, 0, len(messages))
	for _, message := range messages {
		persistedSources = append(persistedSources, message.GetSource())
	}
	if !j.condition.Evaluate(persistedSources) {
		logger.Debug("got message, but condition isn't met yet")

		return nil
	}

	// join sources
	contentBySource := make(map[string]interface{})
	for _, message := range messages {
		var mapTypedContent map[string]interface{}
		err = json.Unmarshal([]byte(message.GetContent()), &mapTypedContent)
		if err == nil {
			logger.Debug(fmt.Sprintf("content of of source %s is a map, putting it as is in the joint message", message.GetSource()))
			contentBySource[message.GetSource()] = mapTypedContent
		} else {
			logger.Debug(fmt.Sprintf("content of of source %s is not a map, putting it as a string in the joint message", message.GetSource()))
			contentBySource[message.GetSource()] = message.GetContent()
		}
	}
	jointContent, err := json.Marshal(contentBySource)
	if err != nil {
		logger.Error("failed to serialize joint message", slog.Any("error", err))

		return err
	}

	jointEvent := event.NewMessage(in.GetKey(), "composed-event", string(jointContent))
	logger.Info("joined message")
	logger.Debug("joint event", slog.String("content", string(jointContent)))

	return next.Call(ctx, jointEvent)
}
