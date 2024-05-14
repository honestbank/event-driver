REPO_NAME=event-driver-extension-google-cloud
docker:
	docker-compose -f integration_test/docker-compose.yaml -p ${REPO_NAME} down;
	docker-compose -f integration_test/docker-compose.yaml -p ${REPO_NAME} up -d;

generate:
	rm -Rf ./mocks
	go generate mockgen.go

test: docker generate
	go test -v -race -covermode=atomic -coverprofile=./cover.out -tags=integration_test ./...
