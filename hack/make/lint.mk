ROOT_DIR := $(shell git -C $(CURDIR) rev-parse --show-toplevel)

.PHONY: lint lint-fix fmt vet

lint:
	cd $(ROOT_DIR) && golangci-lint run

lint-fix:
	cd $(ROOT_DIR) && golangci-lint run --fix

fmt:
	cd $(ROOT_DIR) && go fmt ./...

vet:
	cd $(ROOT_DIR) && go vet ./...
