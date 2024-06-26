package transformer

import (
	"fmt"

	"github.com/honestbank/event-driver/event"
)

// Rule defines a transformer rule that transforms the input event.Message.
// The input event.Message might be updated by the rule.
type Rule func(*event.Message) (*event.Message, error)

func (r Rule) Transform(in *event.Message) (*event.Message, error) {
	return r(in)
}

func (r Rule) append(next Rule) Rule {
	return func(message *event.Message) (*event.Message, error) {
		currentResult, err := r(message)
		if err != nil {
			return nil, err
		}

		return next.Transform(currentResult)
	}
}

// EraseContentFromSources returns a Rule that erases the message content if source is in the list.
func EraseContentFromSources(sources ...string) Rule {
	shouldErase := make(map[string]bool)
	for _, source := range sources {
		shouldErase[source] = true
	}

	return func(message *event.Message) (*event.Message, error) {
		if shouldErase[message.GetSource()] {
			message.SetContent("")
		}

		return message, nil
	}
}

// Identity returns the Rule that keeps the input as is.
func Identity() Rule {
	return func(in *event.Message) (*event.Message, error) {
		return in, nil
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

	return func(message *event.Message) (*event.Message, error) {
		source := message.GetSource()
		if name, isAlias := reverseMap[source]; isAlias {
			message.SetSource(name)
		}

		return message, nil
	}, nil
}
