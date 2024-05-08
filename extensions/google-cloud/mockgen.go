//go:build mockgen

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mock_read_policy.go -package=mocks github.com/lukecold/event-driver/extensions/google-cloud/storage/gcs_event_store ObjectIterator

package main

import (
	_ "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen"
	_ "go.uber.org/mock/mockgen/model"
)
