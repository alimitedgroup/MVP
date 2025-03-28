#!/bin/bash

BASE="http://localhost:8080/api/v1"

echo "Login with global_admin"
TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=global_admin | jq -r '.token')
PARAMS=(-sS -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

echo "User information"
curl "${PARAMS[@]}" -X GET "$BASE/is_logged" | jq
