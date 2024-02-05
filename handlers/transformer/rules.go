package transformer

import (
	"fmt"

	"github.com/lukecold/event-driver/event"
)

// Rule defines a transformer rule that transforms the input event.Message.
// The input event.Message might be updated by the rule.
type Rule func(*event.Message) *event.Message

func (r Rule) Transform(in *event.Message) *event.Message {
	return r(in)
}

func (r Rule) append(next Rule) Rule {
	return func(message *event.Message) *event.Message {
		currentResult := r(message)

		return next.Transform(currentResult)
	}
}

// Identity returns the Rule that keeps the input as is.
func Identity() Rule {
	return func(in *event.Message) *event.Message {
		return in
	}
}

// RenameSources returns a Rule that maps all source aliases to one source name.
// Parameter aliasMap is a map of name to all aliases.
// This function would fail if more than one resource names share the same alias.
func RenameSources(aliasMap map[string][]string) (Rule, error) {
	reverseMap := make(map[string]string)
	for name, aliases := range aliasMap {
		for _, alias := range aliases {
			if _, exist := reverseMap[alias]; exist {
				return nil, fmt.Errorf("alias conflict: %s", alias)
			}

			reverseMap[alias] = name
		}
	}

	return func(message *event.Message) *event.Message {
		source := message.GetSource()
		if name, isAlias := reverseMap[source]; isAlias {
			message.SetSource(name)
		}

		return message
	}, nil
}

// EraseContentFromSources returns a Rule that erases the message content if source is in the list.
func EraseContentFromSources(sources ...string) Rule {
	shouldErase := make(map[string]bool)
	for _, source := range sources {
		shouldErase[source] = true
	}

	return func(message *event.Message) *event.Message {
		if shouldErase[message.GetSource()] {
			message.SetContent("")
		}

		return message
	}
}
