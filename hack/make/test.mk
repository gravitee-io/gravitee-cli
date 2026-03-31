.PHONY: test test-cover test-verbose

test:
	go test ./...

test-cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

test-verbose:
	go test -v ./...
