# Project variables
BINARY_NAME := sshto
OUTPUT_DIR := bin
MODULE := github.com/codoworks/sshto

# Go commands
GO := go
GOFMT := gofmt
GOVET := $(GO) vet
GOTEST := $(GO) test
GOBUILD := $(GO) build

# Build flags
LDFLAGS := -s -w

.PHONY: all build test test-coverage fmt fmt-check vet lint clean install tidy help

## all: Build the binary (default target)
all: build

## build: Build the binary to bin/
build:
	@mkdir -p $(OUTPUT_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) .

## test: Run all tests
test:
	$(GOTEST) -v -race ./...

## test-coverage: Run tests with coverage report
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html

## fmt: Format all Go files
fmt:
	$(GOFMT) -w .

## fmt-check: Check if code is formatted
fmt-check:
	@if [ -n "$$($(GOFMT) -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		$(GOFMT) -d .; \
		exit 1; \
	fi

## vet: Run go vet
vet:
	$(GOVET) ./...

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

## clean: Remove build artifacts
clean:
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.txt coverage.html

## install: Install binary to GOPATH/bin
install:
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(GOPATH)/bin/$(BINARY_NAME) .

## tidy: Tidy and verify module dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
