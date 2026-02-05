# Delete account by ID
# DELETE /api/v1/accounts/:id

curl -s -X DELETE http://localhost:8080/api/v1/accounts/ACCOUNT_UUID_HERE | jq
