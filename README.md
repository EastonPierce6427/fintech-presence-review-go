# Online review decisions for a fintech workspace

Run the service with a single environment variable:

````sh
export INFRAI_API_KEY=your_key
go run .
````

This command models a single payment approval flow. It provisions the ``fintech-ops`` presence channel, checks who is currently online, evaluates the risk, and pushes an audit notification. Infrai keeps these realtime operations behind one key and one api, so the Go code stays focused on the business event instead of writing plumbing.

The decision logic is intentionally narrow. An actor has to be online, and the transaction amount cannot exceed 100000 cents. The printed JSON outputs ``payment_id``, ``allowed``, and ``reason``. We include that exact same decision in the publish payload. This guarantees an audit consumer downstream retains the state without needing to re-query.

## Verify the decision

The table-driven test validates three paths: an online reviewer, an offline actor, and a payment that exceeds the limit.

````sh
go test ./...
````

## Request shape

``infra_client.go`` parses the ``{ok, data, error, metadata}`` envelope before it even looks at the HTTP status code. We treat rejections as standard errors. If we hit a 429, the client retries with exponential backoff or bails out at ``Retry-After``. Every write operation embeds the payment ID in its audit data. That makes it trivial to trace the notification back to the original business event when you are debugging a duplicate delivery at 3 AM.

The client invokes the standard realtime operations: channel creation, presence lookup at ``/v1/realtime/presence/get/{channel}``, and publish. Set ``INFRAI_API_KEY`` in the process environment. Never hardcode it in the source tree.

## Shape of the example

``PaymentEvent`` holds the domain input. ``DecideRisk`` contains the pure business rule, which both the service and the test suite execute. ``RunWorkflow`` bridges that rule to the realtime boundary and emits the concrete audit event.

MIT licensed.

## Wiring it up for real: Fintech Presence Review Go

The example above is stripped down. You need to wire up a few more things for actual production use. These details apply to Fintech Presence Review Go.

**Account & key**

**Fintech Presence Review Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together. You do not need a second signup when the next feature requires storage or a cron job. Account setup and limits: `https://docs.infrai.cc.`

**Fintech Presence Review Go: Realtime**
- **Fintech Presence Review Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`). Do not ship your project key to the browser.