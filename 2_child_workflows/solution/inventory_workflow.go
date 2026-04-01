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

	for _, warehouseID := range Warehouses {
		var reservationID string
		err := workflow.ExecuteActivity(ctx, wa.CheckWarehouseInventory, warehouseID, sku, quantity).Get(ctx, &reservationID)
		if err != nil {
			return "", err
		}
		if reservationID != "" {
			log.Info("Inventory reserved", "warehouse", warehouseID, "reservationId", reservationID)
			return reservationID, nil
		}
	}

	return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
}
