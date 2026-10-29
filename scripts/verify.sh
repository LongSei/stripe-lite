#!/usr/bin/env bash
set -euo pipefail

PORT=8089
SERVER_BIN="./bin/server"
BASE_URL="http://127.0.0.1:${PORT}"

echo "=========================================================="
echo " Starting Stripe-lite API Gateway Self-Verification Suite "
echo "=========================================================="

if [ ! -f "$SERVER_BIN" ]; then
    echo "[-] Binary $SERVER_BIN not found. Building..."
    make build
fi

# Launch server in background with custom port
SERVER_PORT=$PORT $SERVER_BIN &
SERVER_PID=$!

cleanup() {
    echo "[*] Cleaning up server process PID $SERVER_PID..."
    kill -TERM "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
    echo "[+] Server stopped."
}
trap cleanup EXIT

# Wait for server readiness
echo "[*] Polling server readiness on port $PORT..."
READY=0
for i in {1..30}; do
    if curl -s -f -o /dev/null "${BASE_URL}/v1/healthz" 2>/dev/null; then
        READY=1
        break
    fi
    sleep 0.2
done

if [ "$READY" -ne 1 ]; then
    echo "[-] Server failed to become ready within timeout."
    exit 1
fi
echo "[+] Gateway successfully booted and ready."

# Test 1: Liveness / Readiness healthz endpoint
echo -n "[*] Testing GET /v1/healthz ... "
HEALTH_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "${BASE_URL}/v1/healthz")
STATUS_CODE=$(echo "$HEALTH_RESP" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$HEALTH_RESP" | grep -v "HTTP_STATUS")

if [ "$STATUS_CODE" != "200" ]; then
    echo "FAILED (Expected 200, got $STATUS_CODE)"
    exit 1
fi
echo "$BODY" | grep -q '"status":"ok"' || { echo "FAILED (status!=ok)"; exit 1; }
echo "PASSED (HTTP 200, Body: $BODY)"

# Test 2: Root /healthz redirect
echo -n "[*] Testing GET /healthz redirect ... "
REDIRECT_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/healthz")
if [ "$REDIRECT_CODE" != "308" ] && [ "$REDIRECT_CODE" != "301" ] && [ "$REDIRECT_CODE" != "200" ]; then
    echo "FAILED (Unexpected status $REDIRECT_CODE)"
    exit 1
fi
echo "PASSED (HTTP $REDIRECT_CODE)"

# Test 3: Swagger UI HTML endpoint
echo -n "[*] Testing GET /docs (Swagger UI) ... "
DOCS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/docs")
if [ "$DOCS_CODE" != "200" ]; then
    echo "FAILED (Expected 200, got $DOCS_CODE)"
    exit 1
fi
echo "PASSED (HTTP 200)"

# Test 4: OpenAPI JSON endpoint
echo -n "[*] Testing GET /openapi.json ... "
JSON_SPEC=$(curl -s "${BASE_URL}/openapi.json")
[[ "$JSON_SPEC" == *'"title": "Stripe-lite API Gateway"'* ]] || { echo "FAILED (title mismatch)"; exit 1; }
echo "PASSED (Valid JSON specification served)"

# Test 5: OpenAPI YAML endpoint
echo -n "[*] Testing GET /openapi.yaml ... "
YAML_SPEC=$(curl -s "${BASE_URL}/openapi.yaml")
[[ "$YAML_SPEC" == *'title: Stripe-lite API Gateway'* ]] || { echo "FAILED (YAML title mismatch)"; exit 1; }
echo "PASSED (Valid YAML specification served)"

# Test 6: Stub endpoint RFC 7807 Problem Details (POST /v1/payment-intents)
echo -n "[*] Testing Stub POST /v1/payment-intents (RFC 7807 501 Not Implemented) ... "
STUB_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "${BASE_URL}/v1/payment-intents" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: test-key-123" \
    -H "X-Request-Id: 6ba7b810-9dad-11d1-80b4-00c04fd430c8" \
    -d '{"amount": 2000, "currency": "usd", "customer_id": "cus_123"}')
STUB_CODE=$(echo "$STUB_RESP" | grep "HTTP_STATUS" | cut -d':' -f2)
STUB_BODY=$(echo "$STUB_RESP" | grep -v "HTTP_STATUS")

if [ "$STUB_CODE" != "501" ]; then
    echo "FAILED (Expected 501, got $STUB_CODE)"
    exit 1
fi
echo "$STUB_BODY" | grep -q "not-implemented" || { echo "FAILED (RFC 7807 type mismatch)"; exit 1; }
echo "$STUB_BODY" | grep -q "Long" || { echo "FAILED (Domain owner reference missing)"; exit 1; }
echo "PASSED (HTTP 501 RFC 7807 ProblemDocument: $STUB_BODY)"

# Test 7: Stub endpoint RFC 7807 Problem Details (POST /v1/auth/tokens)
echo -n "[*] Testing Stub POST /v1/auth/tokens (RFC 7807 501 Not Implemented) ... "
AUTH_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "${BASE_URL}/v1/auth/tokens" \
    -H "Content-Type: application/json" \
    -d '{"api_key": "sk_test_123"}')
AUTH_CODE=$(echo "$AUTH_RESP" | grep "HTTP_STATUS" | cut -d':' -f2)
AUTH_BODY=$(echo "$AUTH_RESP" | grep -v "HTTP_STATUS")

if [ "$AUTH_CODE" != "501" ]; then
    echo "FAILED (Expected 501, got $AUTH_CODE)"
    exit 1
fi
echo "$AUTH_BODY" | grep -q "Khánh" || { echo "FAILED (Domain owner reference missing)"; exit 1; }
echo "PASSED (HTTP 501 RFC 7807 ProblemDocument)"

# Test 8: Stub endpoint RFC 7807 Problem Details (POST /v1/webhook-endpoints)
echo -n "[*] Testing Stub POST /v1/webhook-endpoints (RFC 7807 501 Not Implemented) ... "
WH_RESP=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "${BASE_URL}/v1/webhook-endpoints" \
    -H "Content-Type: application/json" \
    -d '{"url": "https://example.com/webhook", "events": ["payment_intent.succeeded"]}')
WH_CODE=$(echo "$WH_RESP" | grep "HTTP_STATUS" | cut -d':' -f2)
WH_BODY=$(echo "$WH_RESP" | grep -v "HTTP_STATUS")

if [ "$WH_CODE" != "501" ]; then
    echo "FAILED (Expected 501, got $WH_CODE)"
    exit 1
fi
echo "$WH_BODY" | grep -q "Sơn" || { echo "FAILED (Domain owner reference missing)"; exit 1; }
echo "PASSED (HTTP 501 RFC 7807 ProblemDocument)"

echo "=========================================================="
echo " All Verification Checks Passed Successfully! [8/8]       "
echo "=========================================================="
