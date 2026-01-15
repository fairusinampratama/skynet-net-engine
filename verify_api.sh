#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
KEY="netengine_secret_key_123"
ROUTER_ID=1

echo "🔍 Starting API Verification..."

# 1. Public Health Check
echo -n "1. GET /health ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$HTTP_CODE" -eq 200 ]; then echo "✅ OK"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 2. Router Health (Secured)
echo -n "2. GET /router/$ROUTER_ID/health ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-App-Key: $KEY" "$BASE_URL/router/$ROUTER_ID/health")
if [ "$HTTP_CODE" -eq 200 ]; then echo "✅ OK"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 3. Router Traffic (Secured) - Needs 'user' param
echo -n "3. GET /router/$ROUTER_ID/traffic (Validation) ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-App-Key: $KEY" "$BASE_URL/router/$ROUTER_ID/traffic")
if [ "$HTTP_CODE" -eq 400 ]; then echo "✅ OK (Correctly rejected missing param)"; else echo "❌ FAILED ($HTTP_CODE)"; fi

echo -n "4. GET /router/$ROUTER_ID/traffic?user=UJI_COBA (Mock) ... "
# We expect 500 or timeout if user doesn't exist on router, or 200 if it does. 
# Since we don't know a valid user for sure on the live router, 
# capturing the output might be better, but we check if it hits the API correctly.
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-App-Key: $KEY" "$BASE_URL/router/$ROUTER_ID/traffic?user=UJI_COBA")
if [ "$HTTP_CODE" -ne 404 ] && [ "$HTTP_CODE" -ne 401 ]; then echo "✅ OK (Reachable: $HTTP_CODE)"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 4. Monitoring Feed
echo -n "5. GET /monitoring/targets ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "X-App-Key: $KEY" "$BASE_URL/monitoring/targets")
if [ "$HTTP_CODE" -eq 200 ]; then echo "✅ OK"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 5. Create Secret (Validation)
echo -n "6. POST /secret (Validation) ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "X-App-Key: $KEY" "$BASE_URL/secret")
if [ "$HTTP_CODE" -eq 400 ]; then echo "✅ OK (Correctly rejected empty body)"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 6. Isolate User (Validation)
echo -n "7. POST /isolate (Validation) ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "X-App-Key: $KEY" "$BASE_URL/isolate")
if [ "$HTTP_CODE" -eq 400 ]; then echo "✅ OK (Correctly rejected empty body)"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 7. Sync Router
echo -n "8. POST /sync/$ROUTER_ID ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "X-App-Key: $KEY" "$BASE_URL/sync/$ROUTER_ID")
if [ "$HTTP_CODE" -eq 200 ]; then echo "✅ OK"; else echo "❌ FAILED ($HTTP_CODE)"; fi

# 8. Swagger Public Access
echo -n "9. GET /swagger/index.html ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/index.html")
if [ "$HTTP_CODE" -eq 200 ]; then echo "✅ OK"; else echo "❌ FAILED ($HTTP_CODE)"; fi

echo "🏁 Verification Done."
