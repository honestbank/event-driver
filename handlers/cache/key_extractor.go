package cache

import (
	"github.com/honestbank/event-driver/event"
)

type KeyExtractor interface {
	Extract(*event.Message) (string, error)
}

type getMessageKey struct{}

func (k *getMessageKey) Extract(in *event.Message) (string, error) {
	return in.GetKey(), nil
}

// GetMessageKey returns a CacheKeyExtractor that simply gets the member value "key" from the message.
func GetMessageKey() KeyExtractor {
	return &getMessageKey{}
}
