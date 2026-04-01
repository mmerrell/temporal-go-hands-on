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
// Simulates ~33% of warehouses having stock.
func (a *WarehouseActivities) CheckWarehouseInventory(ctx context.Context, warehouseID, sku string, quantity int) (string, error) {
	activity.GetLogger(ctx).Info("Checking inventory", "warehouse", warehouseID, "sku", sku)
	if warehouseID == "WH-INCHEON" || warehouseID == "WH-BUCHEON" {
		return fmt.Sprintf("RES-%s-%s-%d", warehouseID, sku, time.Now().UnixMilli()), nil
	}
	return "", nil
}
