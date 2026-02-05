#!/usr/bin/env bash
# Stress test: Get account by ID
# 1000 requests, 50 concurrent workers, 5s timeout
#
# Usage: ACCOUNT_ID=<uuid> ./accounts_get.sh

ACCOUNT_ID="${ACCOUNT_ID:?Set ACCOUNT_ID env var}"

hey -n 1000 -c 50 -t 5 \
  http://localhost:8080/api/v1/accounts/${ACCOUNT_ID}
