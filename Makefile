generate:
	rm -Rf ./mocks
	go generate mockgen.go

test: generate
	go test -v -race -coverprofile=./cover.out -covermode=atomic ./...

sync_version:
	for d in ./extensions/*/; do \
 		cd $$d; \
 		go get github.com/lukecold/event-driver && go mod tidy; \
 		cd -; \
	done

commit: test sync_version
	pre-commit install
	pre-commit autoupdate
	pre-commit run -a
