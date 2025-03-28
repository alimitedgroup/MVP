#!/bin/bash

BASE="http://localhost:8080/api/v1"

echo "Login with global_admin"
TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=global_admin | jq -r '.token')
PARAMS=(-sS -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

echo "Create good hat-1"

curl "${PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat, version 1"}' | jq

echo "Get goods status"

curl "${PARAMS[@]}" -X GET "$BASE/goods" | jq

echo "Update good hat-1"

curl "${PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat, version 2"}' | jq

echo "Get goods status"

curl "${PARAMS[@]}" -X GET "$BASE/goods" | jq
