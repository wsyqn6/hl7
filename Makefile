.PHONY: build test test-verbose coverage clean fmt vet tidy help

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=hl7
GO=go
GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt

## build: Build the package
build:
	$(GO) build ./...

## test: Run all tests
test:
	$(GOTEST) ./...

## test-verbose: Run tests with verbose output
test-verbose:
	$(GOTEST) -v ./...

## coverage: Run tests with coverage report
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## bench: Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

## fuzz: Run fuzz tests
fuzz:
	$(GOTEST) -fuzz=Fuzz -fuzztime=10s ./...

## vet: Run go vet
vet:
	$(GOVET) ./...

## fmt: Format code
fmt:
	$(GOFMT) -w .

## tidy: Tidy go modules
tidy:
	$(GOCMD) mod tidy

## clean: Clean build artifacts
clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME) coverage.out coverage.html

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'
