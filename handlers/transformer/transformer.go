package transformer

import (
	"context"
	"log/slog"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
	"github.com/lukecold/event-driver/handlers/options"
)

// transformer implements handlers.Handler that transforms the input with the given rules.
// The input event.Message might be updated by the transformer.
type transformer struct {
	logger *slog.Logger
	rule   Rule
}

func New(rules []Rule, opts ...options.Option) *transformer {
	cfg := options.DefaultOptions()
	for _, opt := range opts {
		opt(&cfg)
	}

	composeRule := Identity()
	for _, rule := range rules {
		composeRule = composeRule.append(rule)
	}

	return &transformer{
		logger: slog.New(slog.NewJSONHandler(cfg.GetLogWriter(), &slog.HandlerOptions{Level: cfg.GetLogLevel()})).
			With(slog.String("handler", "transformer")),
		rule: composeRule,
	}
}

func (m *transformer) WithRules(rules ...Rule) *transformer {
	for _, rule := range rules {
		m.rule = m.rule.append(rule)
	}

	return m
}

func (m *transformer) Process(ctx context.Context, in *event.Message, next handlers.CallNext) error {
	logger := m.logger.With(slog.String("key", in.GetKey()), slog.String("source", in.GetSource()))
	transformed, err := m.rule.Transform(in)
	if err != nil {
		logger.Error("failed to transform message", slog.Any("error", err))

		return err
	}
	logger.Debug("transformed message")

	return next.Call(ctx, transformed)
}
