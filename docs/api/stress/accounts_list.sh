#!/usr/bin/env bash
# Stress test: List accounts
# 1000 requests, 50 concurrent workers, 5s timeout

hey -n 1000 -c 50 -t 5 \
  http://localhost:8080/api/v1/accounts?limit=10&offset=0
