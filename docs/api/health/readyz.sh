# Readiness probe (checks DB connectivity)
# GET /readyz

curl -s http://localhost:8080/readyz | jq
