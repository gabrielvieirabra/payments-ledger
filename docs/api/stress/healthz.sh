#!/usr/bin/env bash
# Stress test: Healthz endpoint
# 5000 requests, 100 concurrent workers, 2s timeout

hey -n 5000 -c 100 -t 2 \
  http://localhost:8080/healthz
