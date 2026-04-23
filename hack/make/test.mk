##@ 🧪 Test

ROOT_DIR ?= $(shell git -C $(CURDIR) rev-parse --show-toplevel)

E2E_COMPOSE := docker compose -f $(ROOT_DIR)/e2e/docker-compose.yml -p gio-e2e

.PHONY: test
test: ## Run unit tests
	@echo "Running unit tests ..."
	@cd $(ROOT_DIR) && go test ./...

.PHONY: test-cover
test-cover: ## Run unit tests with coverage and generate cover.html
	@echo "Running unit tests with coverage ..."
	@cd $(ROOT_DIR) && go test -coverprofile=cover.out ./...
	@cd $(ROOT_DIR) && go tool cover -html=cover.out -o cover.html

.PHONY: e2e-up
e2e-up: ## Start the e2e infra in the background and wait until every service is healthy
	@echo "Starting e2e infra ..."
	@$(E2E_COMPOSE) up -d --wait

.PHONY: e2e-down
e2e-down: ## Stop and remove the e2e infra and its volumes
	@echo "Stopping e2e infra ..."
	@$(E2E_COMPOSE) down -v

.PHONY: test-e2e
test-e2e: ## Run e2e tests against an already-running infra (call e2e-up first)
	@echo "Running e2e tests ..."
	@cd $(ROOT_DIR) && go test -tags e2e -v -timeout 10m ./e2e/

.PHONY: e2e
e2e: ## One-shot e2e cycle: start infra, run tests, always tear down (intended for CI)
	@$(E2E_COMPOSE) up -d --wait
	@cd $(ROOT_DIR) && go test -tags e2e -v -timeout 10m ./e2e/ ; STATUS=$$? ; $(E2E_COMPOSE) down -v ; exit $$STATUS
