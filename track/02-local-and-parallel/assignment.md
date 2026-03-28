---
slug: local-and-parallel
type: challenge
title: 'Exercise 2: Parallel Fan-Out + Local Activities'
teaser: Fan out warehouse checks concurrently and move fast in-process steps to local activities.
notes:
- type: text
  contents: |-
    In Exercise 1, the inventory reservation checked warehouses **one at a time** — sequential
    activity calls. At SailPoint's scale (hundreds of thousands of provisioning workflows per
    month), that's a real latency problem.

    This exercise adds two optimizations:

    **Parallel fan-out:** Check all warehouses simultaneously using `workflow.Go()` goroutines
    and a `workflow.Channel` to collect results. Return as soon as the first warehouse responds
    with stock.

    **Local activities:** `ValidateOrder` and `FraudCheck` are pure in-process checks — no
    network calls. Running them as regular activities wastes two task queue round-trips per
    workflow. Local activities eliminate that overhead and produce a single `MarkerRecorded`
    event instead of the schedule/start/complete triplet.

    Hit **Start** when you're ready.
tabs:
- title: VS Code
  type: service
  hostname: workshop-host
  path: ?folder=/workspace/exercise&openFile=/workspace/exercise/inventory_workflow.go&openFile=/workspace/exercise/fulfillment_workflow.go
  port: 8443
- title: Terminal 1 - Worker
  type: terminal
  hostname: workshop-host
  workdir: /workspace/exercise
- title: Terminal 2 - Starter
  type: terminal
  hostname: workshop-host
  workdir: /workspace/exercise
- title: Temporal Web UI
  type: service
  hostname: workshop-host
  path: /
  port: 8080
difficulty: basic
timelimit: 3600
---

## Exercise 2: Parallel Fan-Out + Local Activities

Open **`inventory_workflow.go`** and **`fulfillment_workflow.go`** in VS Code.
Look for `// TODO` comments — they mark everything you need to implement.

***

### Part A – Create the results channel

In `inventory_workflow.go`, create a channel above the loop to collect results from goroutines:

```go
resultCh := workflow.NewChannel(ctx)
```

***

### Part B – Fan out with `workflow.Go()`

Inside a loop over `Warehouses`, launch one goroutine per warehouse:

```go
for _, warehouseID := range Warehouses {
    wID := warehouseID  // capture loop variable — critical in Go
    workflow.Go(ctx, func(gCtx workflow.Context) {
        var reservationID string
        err := workflow.ExecuteActivity(gCtx, wa.CheckWarehouseInventory, wID, sku, quantity).Get(gCtx, &reservationID)
        if err != nil {
            resultCh.Send(gCtx, "")
            return
        }
        resultCh.Send(gCtx, reservationID)
    })
}
```

> `workflow.Go()` is Temporal's durable goroutine — it replays correctly and is safe
> inside workflow code. Unlike regular Go goroutines, it participates in the workflow sandbox.

***

### Part C – Collect results

Receive one result per goroutine. Return as soon as you find stock:

```go
for i := 0; i < len(Warehouses); i++ {
    var reservationID string
    resultCh.Receive(ctx, &reservationID)
    if reservationID != "" {
        return reservationID, nil
    }
}
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
```

***

### Part D – Add local activities to `FulfillmentWorkflow`

In `fulfillment_workflow.go`, add local activity options and call `ValidateOrder` and
`FraudCheck` before the child workflow:

```go
lao := workflow.LocalActivityOptions{
    StartToCloseTimeout: 5 * time.Second,
}
localCtx := workflow.WithLocalActivityOptions(ctx, lao)
lfa := &LocalFulfillmentActivities{}

if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
    return OrderResult{}, err
}
var riskScore string
if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
    return OrderResult{}, err
}
```

> `workflow.LocalActivityOptions` and `workflow.ExecuteLocalActivity` — not the remote equivalents.
> Everything else (`.Get()`, error handling) is identical.

***

### Part E – Run it

**Terminal 1 - Worker:**
```
go run ./cmd/worker/
```

**Terminal 2 - Starter:**
```
go run ./cmd/starter/
```

In the **Web UI**, open `inventory-ORD-2001`. Notice the `CheckWarehouseInventory` activity
tasks are scheduled with nearly identical timestamps — all running concurrently.

Open `fulfillment-ORD-2001`. Look for `MarkerRecorded` events for the local activities
instead of the schedule/start/complete triplet you'd see for remote activities.

***

Click **Check** when done, or **Solve** to see the reference solution.
