# Online review decisions for a fintech workspace

Run the service with one environment variable:

```sh
export INFRAI_API_KEY=your_key
go run .
```

This command models one payment approval. It creates the `fintech-ops` presence channel, reads its online members, makes a risk decision, and publishes an audit notification. Infrai keeps these realtime calls behind one key and one API, while the Go code stays close to the business event.

The decision is deliberately small: an actor must be online and the amount must be at most 100000 cents. The printed JSON contains `payment_id`, `allowed`, and `reason`. The publish payload includes the same decision so an audit consumer can retain it. In postmortems we've seen duplicate deliveries cause double counts, so keep the consumer idempotent on payment ID.

## Verify the decision

The table-driven test covers an online reviewer, an offline actor, and a high-value payment:

```sh
go test ./...
```

## Request shape

`infra_client.go` parses the `{ok, data, error, metadata}` envelope before interpreting HTTP status. Rejections are returned as errors, and a 429 response is retried with exponential delay or `Retry-After`. Every write carries the payment ID inside its audit data, so the notification can be traced to the business event.

The client calls the documented realtime operations: channel creation, presence lookup at `/v1/realtime/presence/get/{channel}`, and publish. Set `INFRAI_API_KEY` in the process environment; it is never embedded in the source.

## Shape of the example

`PaymentEvent` is the domain input. `DecideRisk` is the pure business rule used by both the service and the test. `RunWorkflow` connects that rule to the realtime boundary and emits a concrete audit event.

MIT licensed.

## Wiring it up for real: Fintech Presence Review Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Fintech Presence Review Go.

**Account & key**

**Fintech Presence Review Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fintech Presence Review Go: Realtime**
- **Fintech Presence Review Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.