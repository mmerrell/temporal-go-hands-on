package fulfillment

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
)

// Warehouses is the list of locations checked during inventory reservation.
var Warehouses = []string{
	"WH-INCHEON", "WH-BUCHEON", "WH-DAEJEON",
	"WH-BUSAN", "WH-GWANGJU", "WH-SEJONG",
}

// WarehouseActivities groups activities used by the inventory reservation child workflow.
type WarehouseActivities struct{}

// CheckWarehouseInventory returns a reservation ID if the warehouse has stock, or "" if not.
func (a *WarehouseActivities) CheckWarehouseInventory(ctx context.Context, warehouseID, sku string, quantity int) (string, error) {
	activity.GetLogger(ctx).Info("Checking inventory", "warehouse", warehouseID, "sku", sku)
	// Simulate: first warehouse that matches SKU prefix wins
	if warehouseID == "WH-INCHEON" || warehouseID == "WH-BUCHEON" {
		return fmt.Sprintf("RES-%s-%s-%d", warehouseID, sku, time.Now().UnixMilli()), nil
	}
	return "", nil
}
