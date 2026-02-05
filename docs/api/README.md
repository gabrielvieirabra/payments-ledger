# API Reference

Base URL: `http://localhost:8080`

## Health

```bash
# Liveness
curl -s http://localhost:8080/healthz | jq

# Readiness (checks DB)
curl -s http://localhost:8080/readyz | jq
```

---

## Accounts

### Create Account
```bash
curl -s -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d @docs/api/accounts/create_account.json | jq
```

Payload (`docs/api/accounts/create_account.json`):
```json
{
  "owner": "John Doe",
  "currency": "BRL"
}
```

Supported currencies: `USD`, `EUR`, `BRL`

### List Accounts
```bash
curl -s "http://localhost:8080/api/v1/accounts?limit=10&offset=0" | jq
```

### Get Account
```bash
curl -s http://localhost:8080/api/v1/accounts/{id} | jq
```

### Delete Account
```bash
curl -s -X DELETE http://localhost:8080/api/v1/accounts/{id} | jq
```

---

## Transactions

### Create Transfer
```bash
curl -s -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d @docs/api/transactions/create_transfer.json | jq
```

Payload (`docs/api/transactions/create_transfer.json`):
```json
{
  "from_account_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "to_account_id": "ffffffff-1111-2222-3333-444444444444",
  "amount": 1500,
  "currency": "BRL"
}
```

> Amount is in the smallest currency unit (e.g. centavos for BRL). Both accounts must share the same currency.

### Get Transaction
```bash
curl -s http://localhost:8080/api/v1/transactions/{id} | jq
```

### List Transactions by Account
```bash
curl -s "http://localhost:8080/api/v1/accounts/{account_id}/transactions?limit=10&offset=0" | jq
```

---

## Entries

### Get Entry
```bash
curl -s http://localhost:8080/api/v1/entries/{id} | jq
```

### List Entries by Account
```bash
curl -s "http://localhost:8080/api/v1/accounts/{account_id}/entries?limit=10&offset=0" | jq
```
