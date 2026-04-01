package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the full order lifecycle using a child workflow
// for inventory reservation.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	fa := &FulfillmentActivities{}

	// Delegate inventory reservation to a child workflow.
	// It gets its own workflow ID, its own event history, and its own retry boundary.
	childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: "inventory-" + order.OrderID,
	})
	var reservationID string
	if err := workflow.ExecuteChildWorkflow(childCtx, InventoryReservationWorkflow,
		order.ItemSKU, order.Quantity).Get(ctx, &reservationID); err != nil {
		return OrderResult{}, err
	}

	var paymentConfirmation string
	if err := workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
		return OrderResult{}, err
	}

	var trackingNumber string
	if err := workflow.ExecuteActivity(ctx, fa.DispatchToFulfillment, order, reservationID).Get(ctx, &trackingNumber); err != nil {
		return OrderResult{}, err
	}

	return OrderResult{
		OrderID:             order.OrderID,
		Status:              "FULFILLED",
		ReservationID:       reservationID,
		PaymentConfirmation: paymentConfirmation,
		TrackingNumber:      trackingNumber,
	}, nil
}
