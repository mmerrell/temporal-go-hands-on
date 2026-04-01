---
slug: parallel-activities
type: challenge
title: 'Exercise 3: Parallel Activities'
teaser: Fan out warehouse checks concurrently using workflow.Go() goroutines and a channel.
notes:
- type: text
  contents: |-
    In Exercise 2, the inventory reservation checked warehouses **one at a time** —
    sequential activity calls. At SailPoint's scale (hundreds of thousands of
    provisioning workflows per month), that is a real latency problem.

    **This exercise makes all warehouse checks run simultaneously.**

    In Go SDK, `workflow.Go()` spawns a durable goroutine inside the workflow sandbox.
    `workflow.NewChannel()` lets goroutines pass results back safely.

    Hit **Start** when you're ready.
tabs:
- title: VS Code
  type: service
  hostname: workshop-host
  path: ?folder=/workspace/exercise&openFile=/workspace/exercise/inventory_workflow.go
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

## Exercise 3: Parallel Activities

Open **`inventory_workflow.go`** in VS Code. Look for `// TODO` comments.

***

### Part A – Create the results channel

Above the loop, create a channel to collect goroutine results:

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

> `workflow.Go()` is Temporal's durable goroutine — safe inside workflow code,
> replays correctly. Unlike regular Go goroutines, it participates in the workflow sandbox.

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

### Part D – Run it

**Terminal 1 - Worker:**
```
go run ./cmd/worker/
```

**Terminal 2 - Starter:**
```
go run ./cmd/starter/
```

In the **Web UI**, open `inventory-ORD-3001`. Notice the `CheckWarehouseInventory`
activity tasks are scheduled with nearly identical timestamps — all running concurrently.
Compare with Exercise 2 where they ran one at a time.

***

Click **Check** when done, or **Solve** to see the reference solution.
