package fulfillment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// FulfillmentActivities groups the three order-processing activities.
// Each simulates a flaky external service call.
type FulfillmentActivities struct{}

// ReserveInventory reserves stock for the order. 30% simulated failure rate.
func (a *FulfillmentActivities) ReserveInventory(ctx context.Context, order Order) (string, error) {
	activity.GetLogger(ctx).Info("Reserving inventory", "orderId", order.OrderID)
	if rand.Float64() < 0.3 {
		return "", temporal.NewApplicationError("inventory service timeout", "InventoryError")
	}
	return fmt.Sprintf("RES-%s-%d", order.ItemSKU, time.Now().UnixMilli()), nil
}

// ProcessPayment charges the customer. 20% simulated failure rate.
func (a *FulfillmentActivities) ProcessPayment(ctx context.Context, order Order) (string, error) {
	activity.GetLogger(ctx).Info("Processing payment", "orderId", order.OrderID)
	if rand.Float64() < 0.2 {
		return "", temporal.NewApplicationError("payment gateway unavailable", "PaymentError")
	}
	return fmt.Sprintf("PAY-%s-%d", order.OrderID, time.Now().UnixMilli()), nil
}

// DispatchToFulfillment ships the order. 20% simulated failure rate.
func (a *FulfillmentActivities) DispatchToFulfillment(ctx context.Context, order Order, reservationID string) (string, error) {
	activity.GetLogger(ctx).Info("Dispatching order", "orderId", order.OrderID)
	if rand.Float64() < 0.2 {
		return "", temporal.NewApplicationError("fulfillment API error", "DispatchError")
	}
	return fmt.Sprintf("TRK-%d-%d", len(reservationID), time.Now().UnixMilli()), nil
}
