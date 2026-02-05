#!/usr/bin/env bash
# Stress test: Create accounts
# 200 requests, 10 concurrent workers, 5s timeout

hey -n 200 -c 10 -t 5 \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{"owner":"stress-test-user","currency":"BRL"}' \
  http://localhost:8080/api/v1/accounts
