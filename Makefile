# Terminal Wrapped - Makefile
# Build cross-platform binaries for distribution

BINARY_NAME=terminal-wrapped
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build targets
.PHONY: all build clean release install test

all: build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Development build with race detector
dev:
	go build -race -o $(BINARY_NAME) .

# Run the binary
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install locally
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# Create release binaries for all platforms
release: clean
	mkdir -p dist
	
	@echo "Building darwin-amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "Building darwin-arm64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Building linux-amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	
	@echo "Building linux-arm64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	
	@echo "Creating checksums..."
	cd dist && shasum -a 256 * > checksums.txt
	
	@echo "Release binaries created in dist/"
	@ls -lah dist/

# Compress release binaries
compress: release
	cd dist && \
	for f in $(BINARY_NAME)-*; do \
		if [ -f "$$f" ] && [ "$$f" != "checksums.txt" ]; then \
			gzip -k "$$f"; \
		fi \
	done
	@echo "Compressed binaries created"
	@ls -lah dist/

# Show binary size
size: build
	@ls -lh $(BINARY_NAME)
	@echo ""
	@file $(BINARY_NAME)

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...

# Update dependencies
deps:
	go mod tidy
	go mod verify

.DEFAULT_GOAL := build

