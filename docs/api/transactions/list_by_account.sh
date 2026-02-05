# List transactions for an account (paginated)
# GET /api/v1/accounts/:account_id/transactions?limit=10&offset=0

curl -s "http://localhost:8080/api/v1/accounts/ACCOUNT_UUID_HERE/transactions?limit=10&offset=0" | jq
