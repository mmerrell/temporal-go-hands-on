package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the order lifecycle using Temporal.
// Temporal handles retries automatically — no manual retry loops needed.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	fa := &FulfillmentActivities{}

	var reservationID string
	if err := workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
		return OrderResult{}, err
	}
	log.Info("Inventory reserved", "reservationId", reservationID)

	var paymentConfirmation string
	if err := workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
		return OrderResult{}, err
	}
	log.Info("Payment confirmed", "confirmation", paymentConfirmation)

	var trackingNumber string
	if err := workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
		return OrderResult{}, err
	}
	log.Info("Order dispatched", "trackingNumber", trackingNumber)

	return OrderResult{
		OrderID:             order.OrderID,
		Status:              "FULFILLED",
		ReservationID:       reservationID,
		PaymentConfirmation: paymentConfirmation,
		TrackingNumber:      trackingNumber,
	}, nil
}
