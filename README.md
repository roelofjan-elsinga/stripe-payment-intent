# Stripe Payment Intent

This Lambda function creates a Stripe Payment Intent.

This needs an environment variable (STRIPE_SECRET) to run.

Input object:

```json
{
  "amount_cents": 999,
  "currency": "eur"
}
```

Output object:

```json
{
  "client_secret": "pi_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```

This needs the API Gateway to proxy an HTTP request to the Lambda function.

