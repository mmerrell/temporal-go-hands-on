package fulfillment

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// InventoryReservationWorkflow fans out to all warehouses in parallel and returns
// the first reservation ID it finds. If no warehouse has stock, it fails with a
// non-retryable error.
func InventoryReservationWorkflow(ctx workflow.Context, sku string, quantity int) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Starting parallel inventory reservation", "sku", sku, "quantity", quantity)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 15 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	wa := &WarehouseActivities{}

	// TODO Part A: Fan out — launch one goroutine per warehouse in parallel.
	//
	// In Go SDK, workflow.Go() spawns a durable goroutine inside the workflow sandbox.
	// Use a workflow.Channel to collect results safely across goroutines.
	//
	// Hint — channel setup:
	//   resultCh := workflow.NewChannel(ctx)
	//
	// Hint — per-warehouse goroutine (call this inside a loop over Warehouses):
	//   wID := warehouseID  // capture loop variable
	//   workflow.Go(ctx, func(gCtx workflow.Context) {
	//       var reservationID string
	//       err := workflow.ExecuteActivity(gCtx, wa.CheckWarehouseInventory, wID, sku, quantity).Get(gCtx, &reservationID)
	//       if err != nil {
	//           resultCh.Send(gCtx, "")
	//           return
	//       }
	//       resultCh.Send(gCtx, reservationID)
	//   })

	// TODO Part B: Collect results — receive one result per warehouse launched.
	//
	// Hint:
	//   for i := 0; i < len(Warehouses); i++ {
	//       var reservationID string
	//       resultCh.Receive(ctx, &reservationID)
	//       if reservationID != "" {
	//           return reservationID, nil
	//       }
	//   }

	// TODO Part C: If no warehouse had stock, return a non-retryable error.
	//
	// Hint:
	//   return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")

	_ = wa        // remove once you use wa above
	_ = temporal.NewApplicationError // remove once you use it above
	return "", nil // replace this
}
