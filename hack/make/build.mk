ROOT_DIR := $(shell git -C $(CURDIR) rev-parse --show-toplevel)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: build build-all clean

build:
	cd $(ROOT_DIR) && go build $(LDFLAGS) -o dist/gio .

build-all:
	cd $(ROOT_DIR) && GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-linux-amd64 .
	cd $(ROOT_DIR) && GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/gio-linux-arm64 .
	cd $(ROOT_DIR) && GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-darwin-amd64 .
	cd $(ROOT_DIR) && GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/gio-darwin-arm64 .
	cd $(ROOT_DIR) && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/gio-windows-amd64.exe .

clean:
	rm -rf $(ROOT_DIR)/dist/
