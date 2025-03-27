TOKEN=$(curl -Ss -X POST localhost:8080/api/v1/login -d username=global_admin | jq -r '.token')

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X PUT localhost:8080/api/v1/goods/hat-1 \
  -d '{"name":"hat","description":"blue hat"}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X PUT localhost:8080/api/v1/goods/hat-2 \
  -d '{"name":"hat","description":"red hat"}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-1/warehouse/1/add \
  -d '{"quantity": 6}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-1/warehouse/2/add \
  -d '{"quantity": 8}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-2/warehouse/1/add \
  -d '{"quantity": 6}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-2/warehouse/2/add \
  -d '{"quantity": 5}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/goods | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/orders \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/orders \
  -d '{"name": "test-order-2", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/orders \
  -d '{"name": "test-order-3", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13, "hat-2": 11}}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/orders | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/goods | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-2/warehouse/2/add \
  -d '{"quantity": 5}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/goods | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/transfers \
  -d '{"receiver_id": "1", "sender_id": "2", "goods": {"hat-2": 2}}' | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/transfers | jq

curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/goods | jq
