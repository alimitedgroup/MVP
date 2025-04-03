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

curl "${GA_PARAMS[@]}" -X POST localhost:8080/api/v1/notifications/queries \
  -d '{"good_id": "hat-1", "operator": "<", "threshold": 100}' | jq

echo "Get notification queries"

curl "${GA_PARAMS[@]}" -X GET localhost:8080/api/v1/notifications/queries | jq