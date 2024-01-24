generate:
	rm -Rf ./mocks
	go generate mockgen.go
