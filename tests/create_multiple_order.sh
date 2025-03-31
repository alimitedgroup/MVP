#!/bin/bash

BASE="http://localhost:8080/api/v1"

echo "Login with global_admin"
GA_TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=global_admin | jq -r '.token')
GA_PARAMS=(-sS -H "Authorization: Bearer $GA_TOKEN" -H "Content-Type: application/json")

echo "Login with local_admin"
LA_TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=local_admin | jq -r '.token')
LA_PARAMS=(-sS -H "Authorization: Bearer $LA_TOKEN" -H "Content-Type: application/json")

echo "Login with client"
C_TOKEN=$(curl -Ss -X POST "$BASE/login" -d username=client | jq -r '.token')
C_PARAMS=(-sS -H "Authorization: Bearer $C_TOKEN" -H "Content-Type: application/json")

echo "Create goods hat-1 and hat-2"

curl "${GA_PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat"}' | jq

curl "${GA_PARAMS[@]}" -X PUT "$BASE/goods/hat-2" \
  -d '{"name":"hat","description":"red hat"}' | jq

echo "Add stock for goods hat-1 and hat-2"

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/2/stock" \
  -d '{"quantity": 8}' | jq

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-2/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-2/warehouse/2/stock" \
  -d '{"quantity": 5}' | jq

echo "Get goods status"

curl "${GA_PARAMS[@]}" -X GET "$BASE/goods" | jq

echo "Create 3 orders"

curl "${C_PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq
curl "${C_PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-2", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq
curl "${C_PARAMS[@]}" -X POST "$BASE/orders" \
  -d '{"name": "test-order-3", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq

echo "Get orders and goods status"

sleep 0.1

curl "${C_PARAMS[@]}" -X GET "$BASE/orders" | jq

curl "${GA_PARAMS[@]}" -X GET "$BASE/goods" | jq
