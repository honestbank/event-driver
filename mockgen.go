//go:build mockgen

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_handlers.go -package=mocks github.com/honestbank/event-driver/handlers CallNext
//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_cache.go -package=mocks github.com/honestbank/event-driver/handlers/cache ConflictResolver,KeyExtractor
//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_event_storage.go -package=mocks github.com/honestbank/event-driver/storage EventStore

package main

import (
	_ "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen"
	_ "go.uber.org/mock/mockgen/model"
)
