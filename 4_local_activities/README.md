# Lab 2: Parallel Activities + Local Activities

In this lab you will make the inventory reservation fan out to all warehouses
simultaneously, and add fast in-process validation steps using local activities.

**What you'll learn:**
- How to run activities in parallel using `workflow.Go()` and `workflow.Channel`
- How local activities differ from remote activities in execution and history
- How to read the Temporal UI to verify both patterns worked correctly

**Time:** ~15 minutes

---

## Background

At the end of Lab 1, the inventory reservation workflow checked warehouses
**one at a time** — sequential activity calls. At SailPoint's scale (hundreds of
thousands of provisioning workflows per month), checking each external system
sequentially is a real latency problem. The fix: fan out in parallel.

Local activities address a different problem: `ValidateOrder` and `FraudCheck`
are pure in-process computations — no network calls. Running them as regular
activities means two unnecessary task queue round-trips to the Temporal Server
per workflow execution. At 20K workflows/day that's 40K wasted Actions. Local
activities eliminate those round-trips and show up in history as a single
`MarkerRecorded` event instead of the schedule/start/complete triplet.

---

## Setup

```bash
temporal server start-dev
```

Open two terminals and `cd` into `lab2-local-and-parallel/practice/`.

```bash
go mod tidy
```

---

## Part A: Fan out with `workflow.Go()`

Open `inventory_workflow.go`.

In Go SDK, `workflow.Go()` spawns a durable goroutine inside the workflow
sandbox. It is **not** a regular Go goroutine — it replays correctly and is safe
to use in workflow code. Use `workflow.NewChannel()` to pass results back to the
main goroutine.

First, create the results channel above your loop:

```go
resultCh := workflow.NewChannel(ctx)
```

Then, inside a loop over `Warehouses`, launch one goroutine per warehouse:

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

> **Go vs. Java:** Java SDK uses `Async.function()` to opt in to concurrency
> (synchronous by default). Go SDK uses `workflow.Go()` goroutines. The mental
> model is the same: fire multiple futures, collect results.

---

## Part B: Collect results

After the launch loop, receive one result per goroutine. Return as soon as you
find stock — don't wait for all warehouses to finish:

```go
for i := 0; i < len(Warehouses); i++ {
    var reservationID string
    resultCh.Receive(ctx, &reservationID)
    if reservationID != "" {
        return reservationID, nil
    }
}
```

> **Go vs. Java:** Java uses `Promise.allOf(promises).get()` to wait for all,
> then iterates. Go's `channel.Receive()` blocks until the next result arrives —
> you can return early as soon as you get a hit without waiting for the rest.

---

## Part C: Handle the out-of-stock case

If all goroutines returned empty strings, no warehouse had stock:

```go
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
```

---

## Part D: Add local activities to `FulfillmentWorkflow`

Open `fulfillment_workflow.go`.

Local activities need their own options type and their own context. Add this
**before** the child workflow call:

```go
lao := workflow.LocalActivityOptions{
    StartToCloseTimeout: 5 * time.Second,
}
localCtx := workflow.WithLocalActivityOptions(ctx, lao)
lfa := &LocalFulfillmentActivities{}

// Validate — returns only error, no result to capture
if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
    return OrderResult{}, err
}

// Fraud check — returns a risk score string
var riskScore string
if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
    return OrderResult{}, err
}
log.Info("Fraud check passed", "riskScore", riskScore)
```

> **Key difference from remote activities:**
> - `workflow.LocalActivityOptions` not `workflow.ActivityOptions`
> - `workflow.WithLocalActivityOptions()` not `workflow.WithActivityOptions()`
> - `workflow.ExecuteLocalActivity()` not `workflow.ExecuteActivity()`
>
> Everything else — `.Get(ctx, &result)`, error handling — is identical.

---

## Part E: Run it

```bash
# Terminal 1
go run ./cmd/worker/

# Terminal 2
go run ./cmd/starter/
```

---

## What to look for in the UI

Open **http://localhost:8233** and find `fulfillment-ORD-2001`.

### Local activities — `FulfillmentWorkflow` history

| Event | What it means |
|---|---|
| `MarkerRecorded` (×2) | Each local activity — no schedule/start/complete, just one marker |
| `ChildWorkflowExecutionStarted` | Inventory child workflow kicked off |
| `ActivityTaskScheduled/Started/Completed` (×2) | ProcessPayment and DispatchToFulfillment — regular remote activities |

The two local activity calls produced **2 events** total. If they were regular
activities they would have produced **6 events** (3 per activity). At 20K
workflows/day, that's 80K fewer Actions.

### Parallel fan-out — `inventory-ORD-2001` history

Open the child workflow. In the event history, look at the
`ActivityTaskScheduled` events for `CheckWarehouseInventory`. In Lab 1 they
appeared one after another (each waited for the previous). Here they should
appear with nearly identical timestamps — all scheduled within the same workflow
task.

| Lab 1 (sequential) | Lab 2 (parallel) |
|---|---|
| Warehouse 1 scheduled → completed → Warehouse 2 scheduled… | All warehouses scheduled in same workflow task |
| Total time ≈ N × single-warehouse latency | Total time ≈ single-warehouse latency |

---

## Solution

The complete solution is in `lab2-local-and-parallel/solution/`. If you want to
compare the parallel channel pattern side by side with the sequential version,
diff `solution/inventory_workflow.go` against Lab 1's `solution/inventory_workflow.go`.
