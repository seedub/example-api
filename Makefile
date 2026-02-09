.PHONY: help build build-amd64 build-arm64 build-all test test-coverage run clean lint install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod verify

build: ## Build the application for current platform
	go build -o bin/example-api ./main.go

build-amd64: ## Build for Linux amd64
	GOOS=linux GOARCH=amd64 go build -o bin/example-api-amd64 ./main.go

build-arm64: ## Build for Linux arm64
	GOOS=linux GOARCH=arm64 go build -o bin/example-api-arm64 ./main.go

build-all: build-amd64 build-arm64 ## Build for all architectures

test: ## Run tests
	go test ./... -v

test-coverage: ## Run tests with coverage report
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML coverage report, run: go tool cover -html=coverage.out"

test-race: ## Run tests with race detection
	go test ./... -race

lint: ## Run linter
	golangci-lint run

run: ## Run the application
	go run main.go

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

.DEFAULT_GOAL := help
