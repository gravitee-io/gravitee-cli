##@ 🔨 Build

## Make this file usable standalone (e.g. IDE gutter run). When included by the
## root Makefile, ROOT_DIR is already set and ?= is a no-op.
ROOT_DIR ?= $(shell git -C $(CURDIR) rev-parse --show-toplevel)

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

## Pull in tool.mk when run standalone so $(GORELEASER) + install-tools resolve.
ifndef GIO_TOOL_MK_LOADED
include $(ROOT_DIR)/hack/make/tool.mk
endif

.PHONY: build
build: ## Build gio for the current platform into dist/gio
	@echo "Building gio $(VERSION) ..."
	@cd $(ROOT_DIR) && go build $(LDFLAGS) -o dist/gio .

.PHONY: release-check
release-check: $(GORELEASER) ## Validate .goreleaser.yaml
	@echo "Checking goreleaser config ..."
	@cd $(ROOT_DIR) && $(GORELEASER) check

.PHONY: release-snapshot
release-snapshot: $(GORELEASER) ## Cross-compile all platforms and build archives into dist/ (no publish)
	@echo "Building snapshot release $(VERSION) ..."
	@cd $(ROOT_DIR) && $(GORELEASER) release --snapshot --clean --skip=publish

## File-target: let any release-* target trigger install-tools when the binary
## is missing.
$(GORELEASER):
	@$(MAKE) install-tools
