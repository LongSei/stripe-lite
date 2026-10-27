"""Check the running gateway against the Payments contract using only the standard library."""
import json
import os
import time
import urllib.error
import urllib.request

base = os.environ["VERIFY_BASE_URL"]
key = os.environ["VERIFY_API_KEY"]
checks = []


def request(method, path, body=None, expected=200, idem=None, authenticated=True):
    headers = {"Content-Type": "application/json"}
    if authenticated:
        headers["X-API-Key"] = key
    if idem is not None:
        headers["Idempotency-Key"] = idem
    data = None if body is None else json.dumps(body).encode()
    req = urllib.request.Request(base + path, data=data, headers=headers, method=method)
    try:
        response = urllib.request.urlopen(req, timeout=5)
    except urllib.error.HTTPError as exc:
        response = exc
    with response:
        status, content_type, raw = response.status, response.headers.get("Content-Type", ""), response.read()
        result = json.loads(raw) if "json" in content_type else raw.decode()
        assert status == expected, (method, path, status, expected, result)
        if expected >= 400:
            assert "application/problem+json" in content_type
            assert result["status"] == expected
        if expected == 401:
            assert response.headers.get("WWW-Authenticate")
        checks.append(f"{method} {path.split('?')[0]} -> {status}")
        return result


for attempt in range(50):
    try:
        health = request("GET", "/v1/healthz")
        assert health["status"] == "ok"
        break
    except (OSError, AssertionError):
        if attempt == 49:
            raise
        time.sleep(0.2)

assert "swagger-ui" in request("GET", "/docs")
spec = request("GET", "/openapi.json")
assert spec["openapi"] == "3.0.3"
assert spec["servers"][0]["url"] == "/v1"
for path, methods in {
    "/payment-intents": ["post", "get"],
    "/payment-intents/{id}": ["get"],
    "/payment-intents/{id}/confirm": ["post"],
    "/payment-intents/{id}/capture": ["post"],
    "/payment-intents/{id}/cancel": ["post"],
    "/refunds": ["post", "get"],
    "/refunds/{id}": ["get"],
}.items():
    assert all(method in spec["paths"][path] for method in methods)
assert "openapi: 3.0.3" in request("GET", "/openapi.yaml")
request("GET", "/v1/payment-intents", expected=401, authenticated=False)
body = dict(amount=2000, currency="usd", customer_id="cus_demo", capture_method="automatic", payment_method_id="pm_demo")
request("POST", "/v1/payment-intents", body, expected=400)
payment = request("POST", "/v1/payment-intents", body, expected=201, idem="create-auto")
pid = payment["id"]
assert request("GET", "/v1/payment-intents/" + pid) == payment
assert request("POST", f"/v1/payment-intents/{pid}/confirm", {})["status"] == "succeeded"
assert request("POST", "/v1/payment-intents", body, expected=201, idem="create-auto") == payment
request("POST", "/v1/payment-intents", dict(body, amount=1000), expected=409, idem="create-auto")
request("POST", f"/v1/payment-intents/{pid}/capture", expected=409)
request("POST", f"/v1/payment-intents/{pid}/cancel", expected=409)
refund_body = dict(payment_intent_id=pid, amount=500, metadata={"order": "example"})
refund = request("POST", "/v1/refunds", refund_body, expected=201, idem="refund-partial")
assert refund["status"] == "pending" and refund["amount"] == 500


def settled(rid):
    for _ in range(50):
        value = request("GET", "/v1/refunds/" + rid)
        if value["status"] != "pending":
            assert value["status"] == "succeeded"
            return value
        time.sleep(0.1)
    raise AssertionError("refund worker did not complete")


settled(refund["id"])
assert request("POST", "/v1/refunds", refund_body, expected=201, idem="refund-partial") == refund
full = request("POST", "/v1/refunds", dict(payment_intent_id=pid), expected=201, idem="refund-full")
assert full["amount"] == 1500
settled(full["id"])
current = request("GET", "/v1/payment-intents/" + pid)
assert current["amount_refunded"] == 2000 and current["amount_refund_pending"] == 0
request("POST", "/v1/refunds", dict(payment_intent_id=pid, amount=1), expected=409, idem="refund-excess")
manual = request("POST", "/v1/payment-intents", dict(body, capture_method="manual"), expected=201, idem="manual")
mid = manual["id"]
assert request("POST", f"/v1/payment-intents/{mid}/confirm", {})["status"] == "requires_capture"
assert request("POST", f"/v1/payment-intents/{mid}/capture")["status"] == "succeeded"
canceled = request("POST", "/v1/payment-intents", body, expected=201, idem="cancel")
cid = canceled["id"]
assert request("POST", f"/v1/payment-intents/{cid}/cancel")["status"] == "canceled"
request("POST", f"/v1/payment-intents/{cid}/confirm", {}, expected=409)
declined = request("POST", "/v1/payment-intents", dict(body, payment_method_id="pm_declined"), expected=201, idem="declined")
request("POST", f"/v1/payment-intents/{declined['id']}/confirm", {}, expected=402)
request("POST", "/v1/payment-intents", dict(body, customer_id="cus_other"), expected=404, idem="other-customer")
page = request("GET", "/v1/payment-intents?limit=1")
assert page["object"] == "list" and len(page["data"]) == 1 and page["has_more"]
request("GET", "/v1/payment-intents?limit=1&starting_after=" + page["data"][0]["id"])
assert len(request("GET", "/v1/refunds?payment_intent_id=" + pid)["data"]) == 2
maximum = request("POST", "/v1/payment-intents", dict(body, amount=9223372036854775807), expected=201, idem="int64")
assert maximum["amount"] == 9223372036854775807
request("POST", "/v1/auth/tokens", dict(api_key="sk_test_sample"), expected=501)
request("POST", "/v1/webhook-endpoints", dict(url="https://example.com/webhook", events=["payment_intent.succeeded"]), expected=501)
print(f"Gateway verification passed: {len(checks)} HTTP assertions, all 9 Payments operations.")
if os.environ.get("VERIFY_REPORT"):
    with open(os.environ["VERIFY_REPORT"], "w") as stream:
        json.dump(dict(passed=True, checks=checks), stream, indent=2)
