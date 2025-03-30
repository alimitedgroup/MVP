# login (valid usernames: global_admin, local_admin, client)
curl -Ss -X POST localhost:8080/api/v1/login \
  -d username=global_admin | jq

# create good (return the good id)
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods \
  -d '{"name":"hat","description":"blue hat"}' | jq

# update a good by its id (or create a good with a specific id)
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X PUT localhost:8080/api/v1/goods/hat-1 \
  -d '{"name":"hat","description":"blue hat"}' | jq

# get all goods
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/goods | jq

# get all warehouses
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/warehouses | jq

# add stock to a good in a warehouse
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/goods/hat-1/warehouse/1/stock \
  -d '{"quantity": 10}' | jq

# remove stock of a good in a warehouse
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X DELETE localhost:8080/api/v1/goods/hat-1/warehouse/1/stock \
  -d '{"quantity": 10}' | jq

# create an order (return the order id)
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/orders \
  -d '{"name": "test-order-1", "full_name": "Mario Rossi", "address": "via roma 12 35012", "goods": {"hat-1": 13}}' | jq

# get all orders
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/orders | jq

# create a transfer (return the transfer id)
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/transfers \
  -d '{"receiver_id": "2", "sender_id": "1", "goods": {"hat-1": 2}}' | jq

# get all transfers
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/transfers | jq

# create a notification query (return the query id)
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X POST localhost:8080/api/v1/notifications/queries \
  -d '{"good_id": "hat-1", "operator": "<", "threshold": "10"}' | jq

# get all notification queries
curl -Ss -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -X GET localhost:8080/api/v1/notifications/queries | jq
