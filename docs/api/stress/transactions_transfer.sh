#!/usr/bin/env bash
# Stress test: Create transfers (concurrent money movement)
# 500 requests, 20 concurrent workers, 10s timeout
#
# Usage:
#   FROM_ACCOUNT_ID=<uuid> TO_ACCOUNT_ID=<uuid> ./transactions_transfer.sh

FROM_ACCOUNT_ID="${FROM_ACCOUNT_ID:?Set FROM_ACCOUNT_ID env var}"
TO_ACCOUNT_ID="${TO_ACCOUNT_ID:?Set TO_ACCOUNT_ID env var}"

hey -n 500 -c 20 -t 10 \
  -m POST \
  -H "Content-Type: application/json" \
  -d "{\"from_account_id\":\"${FROM_ACCOUNT_ID}\",\"to_account_id\":\"${TO_ACCOUNT_ID}\",\"amount\":1,\"currency\":\"BRL\"}" \
  http://localhost:8080/api/v1/transactions
