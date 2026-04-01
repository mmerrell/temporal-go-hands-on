---
slug: child-workflows
type: challenge
title: 'Exercise 2: Child Workflows'
teaser: Extract inventory reservation into a dedicated child workflow with its own history and retry boundary.
notes:
- type: text
  contents: |-
    In the starting code, `FulfillmentWorkflow` calls `CheckWarehouseInventory` as a plain
    activity directly — inventory reservation is just a step in the parent's history.

    **Child workflows give you:**
    - A **separate Event History** — keeps the parent lean and auditable
    - An **independent retry boundary** — the child can fail and retry without restarting the parent
    - A **meaningful workflow ID** you can query independently (e.g., `inventory-ORD-1001`)

    Your job: implement `InventoryReservationWorkflow`, then wire it into `FulfillmentWorkflow`
    as a child workflow call.

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
timelimit: 2400
---

## Exercise 2: Child Workflows

Open **`inventory_workflow.go`** and **`fulfillment_workflow.go`** in VS Code.
Look for `// TODO` comments — they mark everything you need to implement.

Files are in `/workspace/exercise/`.

***

### Part A – Implement `InventoryReservationWorkflow`

In `inventory_workflow.go`, iterate over `Warehouses` and call `CheckWarehouseInventory`
for each warehouse until you find one with stock:

```go
for _, warehouseID := range Warehouses {
    var reservationID string
    err := workflow.ExecuteActivity(ctx, wa.CheckWarehouseInventory, warehouseID, sku, quantity).Get(ctx, &reservationID)
    if err != nil {
        return "", err
    }
    if reservationID != "" {
        return reservationID, nil
    }
}
```

If no warehouse has stock, return a non-retryable error:

```go
return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
```

***

### Part B – Call the child workflow from `FulfillmentWorkflow`

In `fulfillment_workflow.go`, replace the `// REMOVE THIS BLOCK` section with a child workflow call:

```go
childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
    WorkflowID: "inventory-" + order.OrderID,
})
var reservationID string
if err := workflow.ExecuteChildWorkflow(childCtx, InventoryReservationWorkflow,
    order.ItemSKU, order.Quantity).Get(ctx, &reservationID); err != nil {
    return OrderResult{}, err
}
```

> **Note:** `.Get(ctx, &result)` is Go's equivalent of Python's `await` — it blocks
> until the child workflow completes and writes the result into `reservationID`.

***

### Part C – Run it

**Terminal 1 - Worker:**
```
go run ./cmd/worker/
```

**Terminal 2 - Starter:**
```
go run ./cmd/starter/
```

Open the **Temporal Web UI** tab. Find both `fulfillment-ORD-1001` and `inventory-ORD-1001`.
Click into each — notice the child has its **own separate Event History**.

***

Click **Check** when done, or **Solve** to see the reference solution.
