##@ 🧹 Lint

ROOT_DIR ?= $(shell git -C $(CURDIR) rev-parse --show-toplevel)

## Pull in tool.mk when run standalone so $(ADDLICENSE) + install-tools resolve.
## Guard prevents re-inclusion when already loaded by the root Makefile.
ifndef GIO_TOOL_MK_LOADED
include $(ROOT_DIR)/hack/make/tool.mk
endif

## Minimal ignore list: third-party, IDE metadata and build artifacts only.
## Everything else we own must carry a header.
LICENSE_IGNORES := \
	-ignore ".idea/**" \
	-ignore "bin/**" \
	-ignore "dist/**" \
	-ignore "hack/tools/**" \
	-ignore "hack/license.go.txt"

.PHONY: lint
lint: lint-sources lint-licenses ## Run every linter and fail on the first error

.PHONY: lint-sources
lint-sources: ## Run golangci-lint (includes go vet, gofmt, staticcheck, errcheck...)
	@echo "Linting go sources ..."
	@cd $(ROOT_DIR) && golangci-lint run

.PHONY: lint-licenses
lint-licenses: $(ADDLICENSE) ## Check license headers and fail if any file is missing one
	@echo "Checking license headers ..."
	@cd $(ROOT_DIR) && $(ADDLICENSE) -check -f LICENSE_TEMPLATE.txt $(LICENSE_IGNORES) .

.PHONY: add-license
add-license: $(ADDLICENSE) ## Stamp a license header on every file that is missing one
	@echo "Adding license headers ..."
	@cd $(ROOT_DIR) && $(ADDLICENSE) -f LICENSE_TEMPLATE.txt $(LICENSE_IGNORES) .

## File-target for the addlicense binary. When a lint target depends on
## $(ADDLICENSE) and the file is missing, Make runs install-tools to build it.
$(ADDLICENSE):
	@$(MAKE) install-tools
