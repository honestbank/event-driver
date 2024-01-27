package pipeline

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
)

type Pipeline interface {
	WithNextHandler(handler handlers.Handler) Pipeline
	Run(ctx context.Context, in event.Message) error
}

type pipeline struct {
	handlers []handlers.Handler
}

func New() Pipeline {
	return &pipeline{
		handlers: make([]handlers.Handler, 0),
	}
}

func (p *pipeline) WithNextHandler(handler handlers.Handler) Pipeline {
	p.handlers = append(p.handlers, handler)

	return p
}

type next func(ctx context.Context, in event.Message) error

func (n next) Call(ctx context.Context, in event.Message) error {
	return n(ctx, in)
}

func (p *pipeline) Run(pipelineContext context.Context, pipelineInput event.Message) error {
	processor := next(func(ctx context.Context, in event.Message) (_ error) {
		return
	})
	for i := len(p.handlers) - 1; i >= 0; i-- {
		processor = func(ctx context.Context, in event.Message) error {
			return p.handlers[i].Process(ctx, in, processor)
		}
	}

	return processor(pipelineContext, pipelineInput)
}
