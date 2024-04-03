package cache

import (
	"github.com/lukecold/event-driver/event"
)

type KeyExtractor interface {
	Extract(*event.Message) (string, error)
}

type extractMessageKey struct{}

func (k *extractMessageKey) Extract(in *event.Message) (string, error) {
	return in.GetKey(), nil
}

// ExtractMessageKey returns a KeyExtractor that simply gets the member value "key" from the message.
func ExtractMessageKey() KeyExtractor {
	return &extractMessageKey{}
}
