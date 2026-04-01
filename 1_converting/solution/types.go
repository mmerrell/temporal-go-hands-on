package fulfillment

type Order struct {
	OrderID     string  `json:"orderId"`
	CustomerID  string  `json:"customerId"`
	ItemSKU     string  `json:"itemSku"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"totalAmount"`
}

type OrderResult struct {
	OrderID             string `json:"orderId"`
	Status              string `json:"status"`
	ReservationID       string `json:"reservationId"`
	PaymentConfirmation string `json:"paymentConfirmation"`
	TrackingNumber      string `json:"trackingNumber"`
}
