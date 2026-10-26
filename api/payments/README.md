# Payments API contract

The contract defines six PaymentIntent operations and three Refund operations for issue #2.
The URL prefix is `/v1`. The path keys do not include this prefix.

## Files

- `payment_intents.yaml` defines the PaymentIntent operations.
- `refunds.yaml` defines the Refund operations.
- `schemas.yaml` defines the request and response models.

The contract references the shared headers, pagination parameters, and ProblemDocument.
It reuses error responses through references to their definitions in the path files.
The `x-retrieve`, `x-confirm`, `x-capture`, and `x-cancel` fields store reusable path items.
Their references make the operations available at the paths listed below.

## Integration

Keep the existing `/payment-intents` reference in the project root specification.
The owner of that specification must add these entries under `paths`:

```yaml
/payment-intents/{id}:
  $ref: "./payments/payment_intents.yaml#/x-retrieve"
/payment-intents/{id}/confirm:
  $ref: "./payments/payment_intents.yaml#/x-confirm"
/payment-intents/{id}/capture:
  $ref: "./payments/payment_intents.yaml#/x-capture"
/payment-intents/{id}/cancel:
  $ref: "./payments/payment_intents.yaml#/x-cancel"
/refunds:
  $ref: "./payments/refunds.yaml"
/refunds/{id}:
  $ref: "./payments/refunds.yaml#/x-retrieve"
```

The root specification and gateway currently use `/api/v1`.
During integration, align both with the required `/v1` prefix.
Then generate the API code and update the strict handlers.
The new contract replaces the old stub-only 501 response.

## Contract rules

- Creation requires an Idempotency-Key with 1 to 255 characters.
- An identical creation retry returns the original status and response body.
- Automatic confirmation captures the full amount.
- Manual confirmation requires a separate full capture.
- Cancellation is permitted before capture.
- A refund requires a captured payment.
- Pending and successful refunds cannot exceed the captured amount.
- Cursor IDs must match the collection and active filters.
- The default limit is 10. The maximum limit is 100.
- Each operation describes the required permission scope.
- Valid credentials without that permission receive 403.
- List filters with no matches return an empty list.

The YAML descriptions contain the full state rules and examples.
The root specification, generated server, SQL, and runtime handlers are outside this change.
