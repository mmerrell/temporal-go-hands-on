package fulfillment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "fulfillment-tasks"

// FulfillmentWorkflow orchestrates the order lifecycle using Temporal.
func FulfillmentWorkflow(ctx workflow.Context, order Order) (OrderResult, error) {
	log := workflow.GetLogger(ctx)
	log.Info("Processing order", "orderId", order.OrderID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	fa := &FulfillmentActivities{}

	// TODO Part A: Call ReserveInventory and capture the reservationID.
	//
	//   var reservationID string
	//   if err := workflow.ExecuteActivity(ctx, fa.ReserveInventory, order).Get(ctx, &reservationID); err != nil {
	//       return OrderResult{}, err
	//   }

	// TODO Part B: Call ProcessPayment and capture paymentConfirmation.

	// TODO Part C: Call DispatchToFulfillment (passing order AND reservationID)
	//   and capture trackingNumber. Then return a completed OrderResult.

	_ = fa // remove once you use fa above
	return OrderResult{}, nil // replace this
}
