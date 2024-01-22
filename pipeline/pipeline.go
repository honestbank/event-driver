package pipeline

import (
	"context"

	"github.com/KeluDiao/event-driver/event"
	"github.com/KeluDiao/event-driver/handlers"
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

func (p *pipeline) Run(pipelineContext context.Context, pipelineInput event.Message) error {
	processor := func(ctx context.Context, in event.Message) (_ error) {
		return
	}
	for i := len(p.handlers) - 1; i >= 0; i-- {
		processor = func(ctx context.Context, in event.Message) error {
			return p.handlers[i].Process(ctx, in, processor)
		}
	}

	return processor(pipelineContext, pipelineInput)
}