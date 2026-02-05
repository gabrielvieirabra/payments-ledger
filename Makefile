APP_NAME := payments-ledger
BUILD_DIR := bin
MAIN_PATH := ./cmd/api
GO := go

.PHONY: all build run test lint fmt vet clean docker-build docker-up docker-down help \
	stress-health stress-accounts-create stress-accounts-list stress-accounts-get \
	stress-transfers stress-transactions-get stress-entries-list stress-all

all: lint test build

## build: Compile the application binary
build:
	@echo "==> Building $(APP_NAME)..."
	$(GO) build -ldflags="-w -s" -o $(BUILD_DIR)/api $(MAIN_PATH)

## run: Run the application locally (loads .env if present)
run:
	@echo "==> Running $(APP_NAME)..."
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi && $(GO) run $(MAIN_PATH)

## test: Run all tests with race detector
test:
	@echo "==> Running tests..."
	$(GO) test -race -count=1 -coverprofile=coverage.out ./...

## coverage: Show test coverage in browser
coverage: test
	$(GO) tool cover -html=coverage.out -o coverage.html
	open coverage.html

## lint: Run golangci-lint
lint:
	@echo "==> Running linter..."
	golangci-lint run ./...

## fmt: Format source code
fmt:
	@echo "==> Formatting code..."
	$(GO) fmt ./...
	goimports -w .

## vet: Run go vet
vet:
	@echo "==> Running vet..."
	$(GO) vet ./...

## tidy: Tidy go modules
tidy:
	@echo "==> Tidying modules..."
	$(GO) mod tidy

## docker-build: Build Docker image
docker-build:
	@echo "==> Building Docker image..."
	docker build -t $(APP_NAME):latest .

## docker-up: Start all services with docker-compose
docker-up:
	docker compose up -d

## docker-down: Stop all services
docker-down:
	docker compose down

## docker-logs: Tail logs from all services
docker-logs:
	docker compose logs -f

## stress-health: Stress test healthz endpoint (5000 req, 100 concurrent)
stress-health:
	@echo "==> Stress testing /healthz..."
	hey -n 5000 -c 100 -t 2 http://localhost:8080/healthz

## stress-accounts-create: Stress test account creation (200 req, 10 concurrent)
stress-accounts-create:
	@echo "==> Stress testing POST /api/v1/accounts..."
	hey -n 500000 -c 10 -t 5 -m POST \
		-H "Content-Type: application/json" \
		-d '{"owner":"stress-test-user","currency":"BRL"}' \
		http://localhost:8080/api/v1/accounts

## stress-accounts-list: Stress test account listing (1000 req, 50 concurrent)
stress-accounts-list:
	@echo "==> Stress testing GET /api/v1/accounts..."
	hey -n 1000 -c 50 -t 5 "http://localhost:8080/api/v1/accounts?limit=10&offset=0"

## stress-accounts-get: Stress test get account (1000 req, 50 concurrent). Usage: make stress-accounts-get ACCOUNT_ID=<uuid>
stress-accounts-get:
	@test -n "$(ACCOUNT_ID)" || (echo "ERROR: ACCOUNT_ID is required" && exit 1)
	@echo "==> Stress testing GET /api/v1/accounts/$(ACCOUNT_ID)..."
	hey -n 1000 -c 50 -t 5 http://localhost:8080/api/v1/accounts/$(ACCOUNT_ID)

## stress-transfers: Stress test transfers (500 req, 20 concurrent). Usage: make stress-transfers FROM=<uuid> TO=<uuid>
stress-transfers:
	@test -n "$(FROM)" || (echo "ERROR: FROM is required" && exit 1)
	@test -n "$(TO)" || (echo "ERROR: TO is required" && exit 1)
	@echo "==> Stress testing POST /api/v1/transactions..."
	hey -n 500 -c 20 -t 10 -m POST \
		-H "Content-Type: application/json" \
		-d '{"from_account_id":"$(FROM)","to_account_id":"$(TO)","amount":1,"currency":"BRL"}' \
		http://localhost:8080/api/v1/transactions

## stress-transactions-get: Stress test get transaction (1000 req, 50 concurrent). Usage: make stress-transactions-get TRANSACTION_ID=<uuid>
stress-transactions-get:
	@test -n "$(TRANSACTION_ID)" || (echo "ERROR: TRANSACTION_ID is required" && exit 1)
	@echo "==> Stress testing GET /api/v1/transactions/$(TRANSACTION_ID)..."
	hey -n 1000 -c 50 -t 5 http://localhost:8080/api/v1/transactions/$(TRANSACTION_ID)

## stress-entries-list: Stress test list entries (1000 req, 50 concurrent). Usage: make stress-entries-list ACCOUNT_ID=<uuid>
stress-entries-list:
	@test -n "$(ACCOUNT_ID)" || (echo "ERROR: ACCOUNT_ID is required" && exit 1)
	@echo "==> Stress testing GET /api/v1/accounts/$(ACCOUNT_ID)/entries..."
	hey -n 1000 -c 50 -t 5 "http://localhost:8080/api/v1/accounts/$(ACCOUNT_ID)/entries?limit=10&offset=0"

## stress-all: Run all read stress tests sequentially. Usage: make stress-all ACCOUNT_ID=<uuid> TRANSACTION_ID=<uuid>
stress-all: stress-health stress-accounts-list
	@test -n "$(ACCOUNT_ID)" || (echo "WARN: skipping account/entry tests (ACCOUNT_ID not set)" && exit 0)
	$(MAKE) stress-accounts-get ACCOUNT_ID=$(ACCOUNT_ID)
	$(MAKE) stress-entries-list ACCOUNT_ID=$(ACCOUNT_ID)
	@if test -n "$(TRANSACTION_ID)"; then $(MAKE) stress-transactions-get TRANSACTION_ID=$(TRANSACTION_ID); fi

## clean: Remove build artifacts
clean:
	@echo "==> Cleaning..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html tmp

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
