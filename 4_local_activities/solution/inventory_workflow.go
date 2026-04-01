package fulfillment

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// InventoryReservationWorkflow fans out to all warehouses in parallel using workflow.Go()
// goroutines and collects results via a workflow.Channel.
func InventoryReservationWorkflow(ctx workflow.Context, sku string, quantity int) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Starting parallel inventory reservation", "sku", sku, "quantity", quantity)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 15 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	wa := &WarehouseActivities{}

	// Fan out: one goroutine per warehouse, all running concurrently.
	// workflow.Go() is Temporal's durable goroutine — safe to use inside workflow code.
	resultCh := workflow.NewChannel(ctx)

	for _, warehouseID := range Warehouses {
		wID := warehouseID // capture loop variable before goroutine closes over it
		workflow.Go(ctx, func(gCtx workflow.Context) {
			var reservationID string
			err := workflow.ExecuteActivity(gCtx, wa.CheckWarehouseInventory, wID, sku, quantity).Get(gCtx, &reservationID)
			if err != nil {
				log.Warn("Warehouse check failed", "warehouse", wID, "error", err)
				resultCh.Send(gCtx, "")
				return
			}
			resultCh.Send(gCtx, reservationID)
		})
	}

	// Collect one result per warehouse launched.
	// Return as soon as we find stock — don't wait for the rest.
	for i := 0; i < len(Warehouses); i++ {
		var reservationID string
		resultCh.Receive(ctx, &reservationID)
		if reservationID != "" {
			log.Info("Inventory reserved (parallel)", "reservationId", reservationID)
			return reservationID, nil
		}
	}

	return "", temporal.NewApplicationError("no stock available for "+sku, "OutOfStock")
}
