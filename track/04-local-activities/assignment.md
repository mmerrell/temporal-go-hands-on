---
slug: local-activities
type: challenge
title: 'Exercise 4: Local Activities'
teaser: Move fast in-process steps to local activities to reduce task queue round-trips.
notes:
- type: text
  contents: |-
    `ValidateOrder` and `FraudCheck` are pure in-process checks — no network calls,
    no external dependencies. Running them as regular activities wastes two task queue
    round-trips to the Temporal Server per workflow execution.

    At SailPoint's scale (20K+ workflows per day), that is 40K unnecessary Actions.

    **Local activities eliminate the round-trip.** They run inside the worker process
    and produce a single `MarkerRecorded` event in the history instead of the
    schedule/start/complete triplet you get from regular activities.

    Hit **Start** when you're ready.
tabs:
- title: VS Code
  type: service
  hostname: workshop-host
  path: ?folder=/workspace/exercise&openFile=/workspace/exercise/fulfillment_workflow.go&openFile=/workspace/exercise/local_activities.go
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

## Exercise 4: Local Activities

Open **`fulfillment_workflow.go`** and **`local_activities.go`** in VS Code.
The local activity implementations are already in `local_activities.go`.
Your job is to call them from the workflow using the local activity API.

***

### Part A – Set up local activity options

In `fulfillment_workflow.go`, add local activity options **before** the child workflow call:

```go
lao := workflow.LocalActivityOptions{
    StartToCloseTimeout: 5 * time.Second,
}
localCtx := workflow.WithLocalActivityOptions(ctx, lao)
lfa := &LocalFulfillmentActivities{}
```

> Use `workflow.LocalActivityOptions` — not `workflow.ActivityOptions`.
> Use `workflow.WithLocalActivityOptions` — not `workflow.WithActivityOptions`.

***

### Part B – Call `ValidateOrder`

`ValidateOrder` returns only an error — no result to capture:

```go
if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
    return OrderResult{}, err
}
```

***

### Part C – Call `FraudCheck`

`FraudCheck` returns a risk score string:

```go
var riskScore string
if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
    return OrderResult{}, err
}
log.Info("Fraud check passed", "riskScore", riskScore)
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

In the **Web UI**, open `fulfillment-ORD-4001`. Look for `MarkerRecorded` events
for the local activities — compare with the `ActivityTaskScheduled` /
`ActivityTaskStarted` / `ActivityTaskCompleted` triplet for the remote activities.

The two local activity calls produced **2 events** total. As regular activities
they would have produced **6 events**. At 20K workflows/day that is 80K fewer Actions.

***

Click **Check** when done, or **Solve** to see the reference solution.
