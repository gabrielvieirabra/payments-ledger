# List entries for an account (paginated)
# GET /api/v1/accounts/:account_id/entries?limit=10&offset=0

curl -s "http://localhost:8080/api/v1/accounts/ACCOUNT_UUID_HERE/entries?limit=10&offset=0" | jq
