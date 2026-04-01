ROOT_DIR := $(shell git -C $(CURDIR) rev-parse --show-toplevel)

.PHONY: test test-cover test-verbose

test:
	cd $(ROOT_DIR) && go test ./...

test-cover:
	cd $(ROOT_DIR) && go test -coverprofile=cover.out ./...
	cd $(ROOT_DIR) && go tool cover -html=cover.out -o cover.html

test-verbose:
	cd $(ROOT_DIR) && go test -v ./...
