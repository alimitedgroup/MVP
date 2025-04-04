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

echo "Create goods hat-1"

curl "${GA_PARAMS[@]}" -X PUT "$BASE/goods/hat-1" \
  -d '{"name":"hat","description":"blue hat"}' | jq

echo "Add stock for good hat-1"

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/1/stock" \
  -d '{"quantity": 6}' | jq

curl "${LA_PARAMS[@]}" -X POST "$BASE/goods/hat-1/warehouse/2/stock" \
  -d '{"quantity": 2}' | jq

sleep 0.5

echo "Get goods status"

curl "${GA_PARAMS[@]}" -X GET "$BASE/goods" | jq

echo "Create transfer"

curl "${GA_PARAMS[@]}" -X POST "$BASE/transfers" \
  -d '{"receiver_id": "2", "sender_id": "1", "goods": {"hat-1": 5}}' | jq

sleep 0.5

curl "${GA_PARAMS[@]}" -X GET "$BASE/transfers" | jq

curl "${GA_PARAMS[@]}" -X GET "$BASE/goods" | jq
