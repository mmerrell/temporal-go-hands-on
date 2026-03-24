package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the full order lifecycle.
// Lab 2 adds: local activity validation + fraud check before the child workflow.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	// Options for remote activities (task queue round-trip, full retry support)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	fa := &FulfillmentActivities{}

	// TODO Part D: Add local activity validation and fraud check BEFORE the child workflow.
	//
	// Local activities run inside the worker process — no task queue scheduling,
	// no separate history events per call. One MarkerRecorded event per local activity.
	// Use LocalActivityOptions (not ActivityOptions), and workflow.ExecuteLocalActivity
	// (not workflow.ExecuteActivity).
	//
	// Hint — local activity options:
	//   lao := workflow.LocalActivityOptions{
	//       StartToCloseTimeout: 5 * time.Second,
	//   }
	//   localCtx := workflow.WithLocalActivityOptions(ctx, lao)
	//   lfa := &LocalFulfillmentActivities{}
	//
	// Hint — call validate (returns only error):
	//   if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
	//       return OrderResult{}, err
	//   }
	//
	// Hint — call fraud check (returns string result):
	//   var riskScore string
	//   if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
	//       return OrderResult{}, err
	//   }
	//   log.Info("Fraud check passed", "riskScore", riskScore)

	// Child workflow — inventory reservation (fans out in parallel across warehouses)
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
