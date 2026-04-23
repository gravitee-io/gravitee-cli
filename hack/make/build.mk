##@ 🔨 Build

## Make this file usable standalone (e.g. IDE gutter run). When included by the
## root Makefile, ROOT_DIR is already set and ?= is a no-op.
ROOT_DIR ?= $(shell git -C $(CURDIR) rev-parse --show-toplevel)

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: build
build: ## Build gio for the current platform into dist/gio
	@echo "Building gio $(VERSION) ..."
	@cd $(ROOT_DIR) && go build $(LDFLAGS) -o dist/gio .

.PHONY: build-all
build-all: ## Cross-compile gio for linux/macos/windows (amd64 + arm64) into dist/
	@echo "Cross-compiling gio $(VERSION) ..."
	@cd $(ROOT_DIR) && GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-linux-amd64 .
	@cd $(ROOT_DIR) && GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/gio-linux-arm64 .
	@cd $(ROOT_DIR) && GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-darwin-amd64 .
	@cd $(ROOT_DIR) && GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/gio-darwin-arm64 .
	@cd $(ROOT_DIR) && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-windows-amd64.exe .
