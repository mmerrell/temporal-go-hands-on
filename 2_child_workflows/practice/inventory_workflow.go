package fulfillment

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// InventoryReservationWorkflow checks each warehouse in sequence and returns the first
// reservation ID it finds. If no warehouse has stock, it fails with a non-retryable error.
func InventoryReservationWorkflow(ctx workflow.Context, sku string, quantity int) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Starting inventory reservation", "sku", sku, "quantity", quantity)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 15 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	wa := &WarehouseActivities{}

	// TODO Part A: Iterate over Warehouses.
	//   For each warehouseID, call CheckWarehouseInventory and collect the result.
	//   If the result is non-empty, return it immediately (first warehouse with stock wins).
	//
	//   Hint:
	//     var reservationID string
	//     err := workflow.ExecuteActivity(ctx, wa.CheckWarehouseInventory, warehouseID, sku, quantity).Get(ctx, &reservationID)
	//     if err != nil { return "", err }
	//     if reservationID != "" { return reservationID, nil }

	// TODO Part B: If no warehouse had stock, return a non-retryable ApplicationError.
	//
	//   Hint:
	//     return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")

	_ = wa  // remove this line once you use wa above
	_ = temporal.NewApplicationError // remove this line once you use it above
	return "", nil // replace this
}
