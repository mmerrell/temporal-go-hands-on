package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the full order lifecycle:
//  1. Local activity: validate order fields (in-process, no task queue round-trip)
//  2. Local activity: fraud check (in-process)
//  3. Child workflow: inventory reservation (fans out in parallel across warehouses)
//  4. Remote activity: process payment
//  5. Remote activity: dispatch to fulfillment
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	// Remote activity options — full task queue scheduling, durable retries
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Local activity options — runs in-process, no task queue round-trip.
	// Produces a single MarkerRecorded event instead of schedule/start/complete triplet.
	lao := workflow.LocalActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	localCtx := workflow.WithLocalActivityOptions(ctx, lao)

	lfa := &LocalFulfillmentActivities{}
	fa := &FulfillmentActivities{}

	// Step 1: Validate — fast in-process check, no external calls
	if err := workflow.ExecuteLocalActivity(localCtx, lfa.ValidateOrder, order).Get(localCtx, nil); err != nil {
		return OrderResult{}, err
	}

	// Step 2: Fraud check — in-memory rules, no external calls
	var riskScore string
	if err := workflow.ExecuteLocalActivity(localCtx, lfa.FraudCheck, order).Get(localCtx, &riskScore); err != nil {
		return OrderResult{}, err
	}
	log.Info("Fraud check passed", "riskScore", riskScore)

	// Step 3: Reserve inventory via child workflow (parallel fan-out inside)
	childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: "inventory-" + order.OrderID,
	})
	var reservationID string
	if err := workflow.ExecuteChildWorkflow(childCtx, InventoryReservationWorkflow,
		order.ItemSKU, order.Quantity).Get(ctx, &reservationID); err != nil {
		return OrderResult{}, err
	}

	// Step 4: Charge payment
	var paymentConfirmation string
	if err := workflow.ExecuteActivity(ctx, fa.ProcessPayment, order).Get(ctx, &paymentConfirmation); err != nil {
		return OrderResult{}, err
	}

	// Step 5: Dispatch
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
