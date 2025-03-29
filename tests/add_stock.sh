#!/bin/bash

BASE="http://localhost:8080/api/v1"

echo "Login with global_admin"
TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=global_admin | jq -r '.token')
PARAMS=(-sS -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

echo "Create goods hat-1"

curl "${PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat"}' | jq

echo "Add stock for good hat-1"

curl "${PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

echo "Get goods status"

curl "${PARAMS[@]}" -X GET "$BASE/goods" | jq
