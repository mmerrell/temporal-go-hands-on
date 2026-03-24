package fulfillment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// FulfillmentActivities groups the remote activities for the parent workflow.
type FulfillmentActivities struct{}

// ProcessPayment charges the customer. Has a 20% simulated failure rate.
func (a *FulfillmentActivities) ProcessPayment(ctx context.Context, order Order) (string, error) {
	activity.GetLogger(ctx).Info("Processing payment", "orderId", order.OrderID)
	if rand.Float64() < 0.2 {
		return "", temporal.NewApplicationError("payment gateway unavailable", "PaymentError")
	}
	return fmt.Sprintf("PAY-%s-%d", order.OrderID, time.Now().UnixMilli()), nil
}

// DispatchToFulfillment ships the order once inventory is reserved and payment taken.
func (a *FulfillmentActivities) DispatchToFulfillment(ctx context.Context, order Order, reservationID string) (string, error) {
	activity.GetLogger(ctx).Info("Dispatching order", "orderId", order.OrderID)
	if rand.Float64() < 0.2 {
		return "", temporal.NewApplicationError("fulfillment API error", "DispatchError")
	}
	return fmt.Sprintf("TRK-%d-%d", len(reservationID), time.Now().UnixMilli()), nil
}
