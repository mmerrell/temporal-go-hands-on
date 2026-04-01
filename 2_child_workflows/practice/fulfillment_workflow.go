package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the full order lifecycle.
// Currently it calls reserveInventory as a plain activity.
// Your job: replace that with a child workflow.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	fa := &FulfillmentActivities{}

	// TODO Part C: Replace the direct activity call below with a child workflow execution.
	//
	//   1. Set child workflow options with a WorkflowID of "inventory-" + order.OrderID
	//   2. Execute InventoryReservationWorkflow as a child, passing order.ItemSKU and order.Quantity
	//   3. Collect the reservationID result
	//
	//   Hint:
	//     childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
	//         WorkflowID: "inventory-" + order.OrderID,
	//     })
	//     var reservationID string
	//     err := workflow.ExecuteChildWorkflow(childCtx, InventoryReservationWorkflow,
	//         order.ItemSKU, order.Quantity).Get(ctx, &reservationID)
	//     if err != nil { return OrderResult{}, err }

	// REMOVE THIS BLOCK once you implement the child workflow above
	wa := &WarehouseActivities{}
	var reservationID string
	if err := workflow.ExecuteActivity(ctx, wa.CheckWarehouseInventory,
		"WH-INCHEON", order.ItemSKU, order.Quantity).Get(ctx, &reservationID); err != nil {
		return OrderResult{}, err
	}
	// END REMOVE BLOCK

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
