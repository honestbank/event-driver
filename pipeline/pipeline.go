package pipeline

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/utils/reflect"
)

type Pipeline interface {
	// WithNextHandler append a handler to the end of the pipeline.
	WithNextHandler(handler handlers.Handler) Pipeline
	// Process executes the handlers in the pipeline in order.
	// If using customized handlers, please make sure next#Call is executed if it's not the last handler.
	// Process respects context.Deadline, and would return a timeout error immediately when deadline is reached.
	// Please note that the handler may still keep running until it's terminated by itself or
	// when the main goroutine (usually the service) is terminated.
	// One is responsible to make their customized handlers handle timeouts gracefully, to prevent resource leak.
	Process(ctx context.Context, in *event.Message) error
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

type next func(ctx context.Context, in *event.Message) error

func (n next) Call(ctx context.Context, in *event.Message) error {
	return n(ctx, in)
}

func (p *pipeline) Process(ctx context.Context, in *event.Message) error {
	process := next(func(_ context.Context, _ *event.Message) error {
		return nil
	})
	for i := len(p.handlers) - 1; i >= 0; i-- {
		index := i
		processNext := process
		process = func(handlerCtx context.Context, message *event.Message) error {
			errorChan := make(chan error)
			go func(index int, errorChan chan error) {
				errorChan <- p.handlers[index].Process(handlerCtx, message, processNext)
			}(index, errorChan)

			select {
			case gotError := <-errorChan:
				if gotError != nil {
					slog.Error("pipeline failed with error", slog.Int("index", index),
						slog.String("handler", reflect.GetType(p.handlers[index])), slog.Any("error", gotError))
				}

				return gotError
			case <-handlerCtx.Done():
				slog.Error("pipeline timed out", slog.Int("index", index),
					slog.String("handler", reflect.GetType(p.handlers[index])))

				return errors.New("pipeline timed out")
			}
		}
	}

	return process(ctx, in)
}
