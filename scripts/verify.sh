#!/usr/bin/env bash
set -euo pipefail

# Use a disposable database and a restricted login for all smoke tests.
if [[ -z "${PAYMENTS_TEST_DATABASE_URL:-}" ]]; then
    exec ./scripts/with-test-db.sh "$0" "$@"
fi
[[ -x ./bin/server ]] || { echo 'Build the server with make build first.' >&2; exit 1; }
VERIFY_DIR=$(mktemp -d)
SERVER_PID=''
cleanup() {
    if [[ -n "$SERVER_PID" ]]; then kill -TERM "$SERVER_PID" 2>/dev/null || true; wait "$SERVER_PID" 2>/dev/null || true; fi
    rm -rf "$VERIFY_DIR"
}
trap cleanup EXIT
export VERIFY_API_KEY
VERIFY_API_KEY=$(openssl rand -hex 32)
export PAYMENTS_API_KEYS
PAYMENTS_API_KEYS=$(python3 - <<'PY'
import json,os
print(json.dumps([dict(key=os.environ['VERIFY_API_KEY'],merchant_id='merchant_demo',scopes=['payment_intents:read','payment_intents:write','refunds:read','refunds:write'])]))
PY
)
VERIFY_PORT=$(python3 - <<'PY'
import socket
with socket.socket() as s:
    s.bind(('127.0.0.1',0)); print(s.getsockname()[1])
PY
)
export VERIFY_BASE_URL="http://127.0.0.1:$VERIFY_PORT"
APP_ENV=development SERVER_HOST=127.0.0.1 SERVER_PORT="$VERIFY_PORT" SERVER_BASE_PATH=/v1 \
    DATABASE_URL="$PAYMENTS_TEST_DATABASE_URL" PAYMENTS_PROCESSOR=simulated \
    PAYMENTS_DEMO_MERCHANT_ID=merchant_demo PAYMENTS_DEMO_CUSTOMER_ID=cus_demo PAYMENTS_RATE_LIMIT=500 \
    ./bin/server >"$VERIFY_DIR/server.log" 2>&1 &
SERVER_PID=$!
if ! python3 scripts/verify_payments.py; then
    cat "$VERIFY_DIR/server.log" >&2
    exit 1
fi
