package main

import (
	"context"
	"fmt"
	"log"

	"fulfillment"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}
	defer c.Close()

	order := fulfillment.Order{
		OrderID:     "ORD-2001",
		CustomerID:  "CUST-99",
		ItemSKU:     "SKU-WIDGET-XL",
		Quantity:    5,
		TotalAmount: 299.95,
	}

	run, err := c.ExecuteWorkflow(context.Background(),
		client.StartWorkflowOptions{
			ID:        "fulfillment-" + order.OrderID,
			TaskQueue: fulfillment.TaskQueue,
		},
		fulfillment.FulfillmentWorkflow,
		order,
	)
	if err != nil {
		log.Fatalln("Unable to start workflow:", err)
	}

	var result fulfillment.OrderResult
	if err := run.Get(context.Background(), &result); err != nil {
		log.Fatalln("Workflow failed:", err)
	}

	fmt.Printf("Order complete: %+v\n", result)
}
