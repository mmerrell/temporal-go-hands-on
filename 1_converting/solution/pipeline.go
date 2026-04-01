package fulfillment

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// FulfillmentPipeline is the "before Temporal" version of order fulfillment.
//
// Problems to spot:
//   - Manual retry loops with time.Sleep — not durable, lost on crash
//   - State in local variables — if the process dies after payment, there is no record
//   - No visibility into which step is running
//   - Double-charge risk: payment succeeded but dispatch failed, caller retries from scratch

const (
	maxRetries   = 5
	retryDelayMs = 2000
)

func RunPipeline(order Order) (OrderResult, error) {
	log.Printf("Starting fulfillment for order %s", order.OrderID)

	var reservationID, paymentConfirmation, trackingNumber string
	var err error

	// Step 1: Reserve inventory — manual retry loop
	for attempt := range maxRetries {
		reservationID, err = reserveInventory(order)
		if err == nil {
			break
		}
		if attempt == maxRetries-1 {
			return OrderResult{}, fmt.Errorf("reserveInventory failed after %d attempts: %w", maxRetries, err)
		}
		log.Printf("reserveInventory attempt %d failed: %v — retrying in %dms", attempt+1, err, retryDelayMs)
		time.Sleep(retryDelayMs * time.Millisecond)
	}
	log.Printf("Inventory reserved: %s", reservationID)

	// Step 2: Process payment — manual retry loop
	for attempt := range maxRetries {
		paymentConfirmation, err = processPayment(order)
		if err == nil {
			break
		}
		if attempt == maxRetries-1 {
			return OrderResult{}, fmt.Errorf("processPayment failed after %d attempts: %w", maxRetries, err)
		}
		log.Printf("processPayment attempt %d failed: %v — retrying in %dms", attempt+1, err, retryDelayMs)
		time.Sleep(retryDelayMs * time.Millisecond)
	}
	log.Printf("Payment confirmed: %s", paymentConfirmation)

	// Step 3: Dispatch — manual retry loop
	for attempt := range maxRetries {
		trackingNumber, err = dispatchToFulfillment(order, reservationID)
		if err == nil {
			break
		}
		if attempt == maxRetries-1 {
			return OrderResult{}, fmt.Errorf("dispatchToFulfillment failed after %d attempts: %w", maxRetries, err)
		}
		log.Printf("dispatchToFulfillment attempt %d failed: %v — retrying in %dms", attempt+1, err, retryDelayMs)
		time.Sleep(retryDelayMs * time.Millisecond)
	}
	log.Printf("Dispatched, tracking: %s", trackingNumber)

	return OrderResult{
		OrderID:             order.OrderID,
		Status:              "FULFILLED",
		ReservationID:       reservationID,
		PaymentConfirmation: paymentConfirmation,
		TrackingNumber:      trackingNumber,
	}, nil
}

// ── Simulated service calls ───────────────────────────────────────────────────

func reserveInventory(order Order) (string, error) {
	if rand.Float64() < 0.3 {
		return "", fmt.Errorf("inventory service timeout")
	}
	return fmt.Sprintf("RES-%s-%d", order.ItemSKU, time.Now().UnixMilli()), nil
}

func processPayment(order Order) (string, error) {
	if rand.Float64() < 0.2 {
		return "", fmt.Errorf("payment gateway unavailable")
	}
	return fmt.Sprintf("PAY-%s-%d", order.OrderID, time.Now().UnixMilli()), nil
}

func dispatchToFulfillment(order Order, reservationID string) (string, error) {
	if rand.Float64() < 0.2 {
		return "", fmt.Errorf("fulfillment API error")
	}
	return fmt.Sprintf("TRK-%d-%d", len(reservationID), time.Now().UnixMilli()), nil
}
