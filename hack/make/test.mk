ROOT_DIR := $(shell git -C $(CURDIR) rev-parse --show-toplevel)
E2E_COMPOSE := docker compose -f $(ROOT_DIR)/e2e/docker-compose.yml -p gio-e2e

.PHONY: test test-cover test-verbose test-e2e e2e e2e-up e2e-down

test:
	cd $(ROOT_DIR) && go test ./...

test-cover:
	cd $(ROOT_DIR) && go test -coverprofile=cover.out ./...
	cd $(ROOT_DIR) && go tool cover -html=cover.out -o cover.html

test-verbose:
	cd $(ROOT_DIR) && go test -v ./...

## e2e-up: start the e2e infra in the background and wait until healthy.
e2e-up:
	$(E2E_COMPOSE) up -d --wait

## e2e-down: stop and remove the e2e infra and its volumes.
e2e-down:
	$(E2E_COMPOSE) down -v

## test-e2e: run e2e tests. Assumes infra is already up - run 'make e2e-up' first.
test-e2e:
	cd $(ROOT_DIR) && go test -tags e2e -v -timeout 10m ./e2e/

## e2e: one-shot cycle - start infra, run tests, always tear down. Intended for CI.
e2e:
	$(E2E_COMPOSE) up -d --wait
	@cd $(ROOT_DIR) && go test -tags e2e -v -timeout 10m ./e2e/ ; STATUS=$$? ; $(E2E_COMPOSE) down -v ; exit $$STATUS
