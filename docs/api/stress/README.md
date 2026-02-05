# Stress Tests (hey)

Load tests using [hey](https://github.com/rakyll/hey).

## Install

```bash
# macOS
brew install hey

# Go
go install github.com/rakyll/hey@latest
```

## Setup

Before running transfer tests, create two accounts and seed balance:

```bash
# Create source account
curl -s -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{"owner":"Alice","currency":"BRL"}' | jq

# Create destination account
curl -s -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{"owner":"Bob","currency":"BRL"}' | jq
```

## Run

```bash
# Health endpoint (baseline)
./stress/healthz.sh

# Account creation
./stress/accounts_create.sh

# Account listing (read)
./stress/accounts_list.sh

# Account get (read, single)
ACCOUNT_ID=<uuid> ./stress/accounts_get.sh

# Transfer (write, concurrent, tests locking)
FROM_ACCOUNT_ID=<uuid> TO_ACCOUNT_ID=<uuid> ./stress/transactions_transfer.sh

# Transaction get (read)
TRANSACTION_ID=<uuid> ./stress/transactions_get.sh

# Entries list (read)
ACCOUNT_ID=<uuid> ./stress/entries_list.sh
```

## What to look for

- **Latency distribution** — p50, p95, p99
- **Error rate** — any non-2xx responses under load
- **Throughput** — requests/sec
- **Transfer consistency** — after `transactions_transfer.sh`, verify `from_account.balance + to_account.balance == original_total`
