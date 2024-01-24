package handlers

import (
	"context"

	"github.com/KeluDiao/event-driver/event"
)

type CallNext interface {
	Call(ctx context.Context, in event.Message) error
}

type Handler interface {
	Process(ctx context.Context, in event.Message, next CallNext) error
}
