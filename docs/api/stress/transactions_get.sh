#!/usr/bin/env bash
# Stress test: Get transaction by ID
# 1000 requests, 50 concurrent workers, 5s timeout
#
# Usage: TRANSACTION_ID=<uuid> ./transactions_get.sh

TRANSACTION_ID="${TRANSACTION_ID:?Set TRANSACTION_ID env var}"

hey -n 1000 -c 50 -t 5 \
  http://localhost:8080/api/v1/transactions/${TRANSACTION_ID}
