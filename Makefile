# terminal-wrapped - makefile
# build cross-platform binaries for distribution

BINARY_NAME=terminal-wrapped
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# build targets
.PHONY: all build clean release install test

all: build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# development build with race detector
dev:
	go build -race -o $(BINARY_NAME) .

# run the binary
run: build
	./$(BINARY_NAME)

# run tests
test:
	go test -v ./...

# clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# install locally
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# create release binaries for all platforms
release: clean
	mkdir -p dist
	
	@echo "building darwin-amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "building darwin-arm64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "building linux-amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	
	@echo "building linux-arm64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	
	@echo "creating checksums..."
	cd dist && shasum -a 256 * > checksums.txt
	
	@echo "release binaries created in dist/"
	@ls -lah dist/

# compress release binaries
compress: release
	cd dist && \
	for f in $(BINARY_NAME)-*; do \
		if [ -f "$$f" ] && [ "$$f" != "checksums.txt" ]; then \
			gzip -k "$$f"; \
		fi \
	done
	@echo "Compressed binaries created"
	@ls -lah dist/

# show binary size
size: build
	@ls -lh $(BINARY_NAME)
	@echo ""
	@file $(BINARY_NAME)

# format code
fmt:
	go fmt ./...

# lint code
lint:
	go vet ./...

# update dependencies
deps:
	go mod tidy
	go mod verify

.DEFAULT_GOAL := build

