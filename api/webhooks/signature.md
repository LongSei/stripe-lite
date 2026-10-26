# Webhook signature verification

Webhook deliveries include a `Stripe-Signature` header so the receiver can
verify that the request was signed with the endpoint's signing secret. The
signing secret has the `whsec_` prefix and must be kept confidential.

## Header format

The header uses this format:

```text
Stripe-Signature: t=timestamp,v1=signature
```

- `t` is the Unix timestamp in seconds when the request was signed.
- `v1` is the lowercase hexadecimal HMAC-SHA256 digest.

## Signing

The sender builds the signed payload by concatenating the timestamp, a literal
period, and the exact raw request body:

```text
signed_payload = timestamp + "." + raw_body
```

The signature is calculated as:

```text
signature = HMAC-SHA256(secret, signed_payload)
```

The secret is the endpoint signing secret with the `whsec_` prefix. The
signature is encoded as lowercase hexadecimal in the `v1` header value.

## Verification

1. Parse the `t` and `v1` values from `Stripe-Signature`.
2. Read and preserve the raw request body before JSON parsing or reformatting.
3. Reconstruct `signed_payload` as `t + "." + raw_body`.
4. Compute HMAC-SHA256 with the endpoint secret.
5. Compare the computed lowercase hexadecimal digest with `v1` using a
   constant-time comparison.
6. Reject requests whose timestamp is outside an application-defined tolerance
   to reduce replay risk.
7. Parse the JSON body only after signature verification succeeds.

The receiver must use the secret currently valid for the endpoint. After a
successful secret rotation, the old secret is invalidated and the new secret
must be used for subsequent deliveries.

## Example delivery

```http
POST /api/webhooks HTTP/1.1
Host: example.com
Content-Type: application/json
Stripe-Signature: t=1792945688,v1=4f5c2e0e5f7c4c4a4b8c5a6e2d1f9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f

{"id":"evt_01HXZ5N7Q8R901","object":"event","type":"payment_intent.succeeded","created":1792945688,"data":{"id":"pi_01HXZ5JABC1234","amount":2000,"currency":"usd","status":"succeeded","created_at":"2026-10-25T16:28:08Z"}}
```

The example digest is illustrative. Implementations must calculate `v1` from
the actual secret, timestamp, and exact raw body bytes.
