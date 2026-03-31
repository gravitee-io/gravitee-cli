.PHONY: lint lint-fix fmt vet

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

fmt:
	go fmt ./...

vet:
	go vet ./...
