generate:
	rm -Rf ./mocks
	go generate mockgen.go

test: generate
	go test -v -race -coverprofile=./cover.out -covermode=atomic ./...
