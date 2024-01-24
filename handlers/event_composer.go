package handlers

import (
	"context"
	"encoding/json"

	"github.com/KeluDiao/event-driver/event"
	"github.com/KeluDiao/event-driver/storage"
)

// TODO: implement me
// requirement: be able to match (AND|OR|()) operations
type SourceMatcher interface {
	Match(sources []string) bool
}

// EventComposer composes the events with the same key when the sources match the criteria given by SourceMatcher
type EventComposer struct {
	storage       storage.EventStore
	sourceMatcher SourceMatcher
}

func (e *EventComposer) Process(ctx context.Context, in event.Message, next CallNext) error {
	// persist input message by key & source
	err := e.storage.Persist(in.GetKey(), in.GetSource(), in.GetContent())
	if err != nil {
		return err
	}

	// validate sources
	messages, err := e.storage.LookUpByKey(in.GetKey())
	if err != nil {
		return err
	}
	persistedSources := make([]string, 0, len(messages))
	for _, message := range messages {
		persistedSources = append(persistedSources, message.GetSource())
	}
	if !e.sourceMatcher.Match(persistedSources) {
		return nil
	}

	// compose
	contentBySource := make(map[string]string)
	for _, message := range messages {
		contentBySource[message.GetSource()] = message.GetContent()
	}
	composedContent, err := json.Marshal(contentBySource)
	if err != nil {
		return err
	}

	composedInput := event.NewMessage(in.GetKey(), "composed-event", string(composedContent))

	return next.Call(ctx, composedInput)
}
