//go:build mockgen

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_handlers.go -package=mocks github.com/KeluDiao/event-driver/handlers CacheHitResolver,CallNext
//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_event_storage.go -package=mocks github.com/KeluDiao/event-driver/storage EventStore

package main

import (
	_ "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen"
	_ "go.uber.org/mock/mockgen/model"
)
