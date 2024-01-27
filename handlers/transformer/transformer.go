package transformer

import (
	"context"

	"github.com/lukecold/event-driver/event"
	"github.com/lukecold/event-driver/handlers"
)

// transformer transforms the input with the given rules
type transformer struct {
	rule Rule
}

func New(rules ...Rule) *transformer {
	composeRule := Identity()
	for _, rule := range rules {
		composeRule = composeRule.append(rule)
	}

	return &transformer{
		rule: composeRule,
	}
}

func (m *transformer) WithRules(rules ...Rule) *transformer {
	for _, rule := range rules {
		m.rule = m.rule.append(rule)
	}

	return m
}

func (m *transformer) Process(ctx context.Context, in event.Message, next handlers.CallNext) error {
	transformed := m.rule.Transform(in)

	return next.Call(ctx, transformed)
}
