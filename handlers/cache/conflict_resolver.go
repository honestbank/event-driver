package cache

import (
	"context"

	"github.com/honestbank/event-driver/event"
	"github.com/honestbank/event-driver/handlers"
)

// ConflictResolver implements handlers.Handler that resolves the case
// when the input matches an existing record in the cache.
type ConflictResolver interface {
	Resolve(ctx context.Context, in *event.Message, next handlers.CallNext) error
}

type skipOnConflict struct{}

func (s *skipOnConflict) Resolve(_ context.Context, _ *event.Message, _ handlers.CallNext) error {
	return nil
}

// SkipOnConflict returns a ConflictResolver that simply skips the following handlers in pipeline on conflict.
func SkipOnConflict() ConflictResolver {
	return &skipOnConflict{}
}
