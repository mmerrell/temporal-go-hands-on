package main

import (
	"log"

	"fulfillment"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}
	defer c.Close()

	w := worker.New(c, fulfillment.TaskQueue, worker.Options{})

	w.RegisterWorkflow(fulfillment.FulfillmentWorkflow)
	w.RegisterWorkflow(fulfillment.InventoryReservationWorkflow)
	w.RegisterActivity(&fulfillment.FulfillmentActivities{})
	w.RegisterActivity(&fulfillment.WarehouseActivities{})

	log.Println("Worker started on task queue:", fulfillment.TaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker error:", err)
	}
}
