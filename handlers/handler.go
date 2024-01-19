package handlers

import (
	"context"

	"github.com/KeluDiao/event-driver/event"
)

type Handler interface {
	Process(ctx context.Context, in event.Message, next func(ctx context.Context, in event.Message) error) error
}
