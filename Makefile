generate:
	rm -Rf ./mocks
	go generate mockgen.go

test: generate
	go test -v -race -coverprofile=./cover.out -covermode=atomic ./...

commit: test
	pre-commit install
	pre-commit autoupdate
	pre-commit run -a
