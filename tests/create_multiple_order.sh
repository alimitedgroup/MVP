#!/bin/bash

BASE="http://localhost:8080/api/v1"

echo "Login with global_admin"
TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=global_admin | jq -r '.token')
PARAMS=(-sS -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

echo "Create goods hat-1 and hat-2"

curl "${PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat"}' | jq

curl "${PARAMS[@]}" -X PUT "$BASE/goods/hat-2" \
  -d '{"name":"hat","description":"red hat"}' | jq

echo "Add stock for goods hat-1 and hat-2"

curl "${PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

curl "${PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/2/stock" \
  -d '{"quantity": 8}' | jq

curl "${PARAMS[@]}" -X POST "$BASE/goods/hat-2/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

curl "${PARAMS[@]}" -X POST "$BASE/goods/hat-2/warehouse/2/stock" \
  -d '{"quantity": 5}' | jq

echo "Get goods status"

curl "${PARAMS[@]}" -X GET "$BASE/goods" | jq

echo "Create 3 orders"

curl "${PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq
curl "${PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq
curl "${PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq

echo "Get orders and goods status"

curl "${PARAMS[@]}" -X GET "$BASE/orders" | jq

curl "${PARAMS[@]}" -X GET "$BASE/goods" | jq
