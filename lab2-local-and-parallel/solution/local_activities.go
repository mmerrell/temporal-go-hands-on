package fulfillment

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/temporal"
)

// LocalFulfillmentActivities contains fast, in-process checks that don't make
// external network calls. They run inside the worker process — no task queue
// round-trip, no separate schedule/start/complete history events.
type LocalFulfillmentActivities struct{}

// ValidateOrder checks that the order fields are well-formed.
// Fast, in-process — ideal as a local activity.
func (a *LocalFulfillmentActivities) ValidateOrder(_ context.Context, order Order) error {
	if order.OrderID == "" {
		return temporal.NewApplicationError("order ID is required", "ValidationError")
	}
	if order.ItemSKU == "" {
		return temporal.NewApplicationError("item SKU is required", "ValidationError")
	}
	if order.Quantity <= 0 {
		return temporal.NewApplicationError(
			fmt.Sprintf("quantity must be positive, got %d", order.Quantity), "ValidationError")
	}
	return nil
}

// FraudCheck runs an in-memory fraud rule against the order.
// Returns a risk score string. Fast, in-process — ideal as a local activity.
func (a *LocalFulfillmentActivities) FraudCheck(_ context.Context, order Order) (string, error) {
	if order.TotalAmount > 10000 {
		return "", temporal.NewApplicationError(
			fmt.Sprintf("order %s flagged: amount %.2f exceeds fraud threshold", order.OrderID, order.TotalAmount),
			"FraudDetected",
		)
	}
	return "CLEAR", nil
}
