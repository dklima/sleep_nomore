.PHONY: all build test clean install lint

BINARY_NAME=sleepnomore
BINARY_AMD64=$(BINARY_NAME)_amd64.exe
BINARY_ARM64=$(BINARY_NAME)_arm64.exe

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean

LDFLAGS=-ldflags="-H windowsgui -s -w"

all: test build

build: build-amd64 build-arm64

build-amd64:
	@echo "Building for Windows x86_64..."
	@set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BINARY_AMD64) -v

build-arm64:
	@echo "Building for Windows ARM64..."
	@set GOOS=windows&& set GOARCH=arm64&& $(GOBUILD) $(LDFLAGS) -o $(BINARY_ARM64) -v

test:
	@echo "Running tests..."
	@$(GOTEST) -v -short ./...

test-all:
	@echo "Running all tests including integration..."
	@$(GOTEST) -v ./...
	@$(GOTEST) -v -tags=integration ./...

bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

coverage:
	@echo "Generating test coverage..."
	@$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@if exist $(BINARY_AMD64) del $(BINARY_AMD64)
	@if exist $(BINARY_ARM64) del $(BINARY_ARM64)
	@if exist coverage.out del coverage.out
	@if exist coverage.html del coverage.html

deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

lint:
	@echo "Running linter..."
	@golangci-lint run --enable-all

install: build
	@echo "Installing $(BINARY_AMD64)..."
	@copy $(BINARY_AMD64) "%USERPROFILE%\AppData\Local\Microsoft\WindowsApps\"

help:
	@echo "Available targets:"
	@echo "  make build      - Build for all architectures"
	@echo "  make test       - Run unit tests"
	@echo "  make test-all   - Run all tests including integration"
	@echo "  make bench      - Run benchmarks"
	@echo "  make coverage   - Generate test coverage report"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make lint       - Run linter (requires golangci-lint)"
	@echo "  make install    - Install to user's PATH"
	@echo "  make help       - Show this help message"