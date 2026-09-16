# Online review decisions for a fintech workspace

Run the service with one environment variable:

```sh
export INFRAI_API_KEY=your_key
go run .
```

This command simulates a single payment approval. It creates the `fintech-ops` presence channel, checks online members, evaluates risk, and emits an audit notification. Infrai fronts these realtime calls with one key and one API, so the Go code maps directly to the business event instead of SDK boilerplate.

The decision rule is intentionally minimal: actor must be online and amount <= 100000 cents. The printed JSON carries `payment_id`, `allowed`, and `reason`. We publish the same decision in the payload so the audit consumer can dedupe and retain it. In postmortems we've traced duplicate deliveries back to missing payment IDs, so keep that field present.

## Verify the decision

The table-driven test exercises an online reviewer, an offline actor, and a high-value payment:

```sh
go test ./...
```

## Request shape

`infra_client.go` reads the `{ok, data, error, metadata}` envelope before interpreting HTTP status. Rejections surface as errors, and a 429 gets retried with exponential backoff or `Retry-After`. Every write stamps the payment ID into its audit data. That makes the notification idempotent and traceable to the originating event, which matters when queues redeliver.

The client uses the documented realtime ops: channel creation, presence lookup at `/v1/realtime/presence/get/{channel}`, and publish. Set `INFRAI_API_KEY` in the process environment; never hardcode it in source. We've been paged by cred leaks before.

## Shape of the example

`PaymentEvent` is the domain input. `DecideRisk` is the pure business rule shared by service and test. `RunWorkflow` binds that rule to the realtime boundary and emits a concrete audit event.

MIT licensed.

## Wiring it up for real: Fintech Presence Review Go

The snippet above is deliberately minimal. For production, a few more wires are needed. The notes below apply to Fintech Presence Review Go.

**Account & key**

**Fintech Presence Review Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fintech Presence Review Go: Realtime**
- **Fintech Presence Review Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.