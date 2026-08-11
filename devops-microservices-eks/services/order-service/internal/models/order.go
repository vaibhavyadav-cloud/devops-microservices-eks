package models

import "time"

type Order struct {
	ID            string    `json:"id"`
	CustomerEmail string    `json:"customerEmail"`
	Item          string    `json:"item"`
	Quantity      int       `json:"quantity"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}

// `binding` tags are Gin's validation rules (via go-playground/validator
// under the hood) — the Go equivalent of jakarta.validation annotations in
// the Notification service and express-validator-style checks in Auth.
type CreateOrderRequest struct {
	CustomerEmail string `json:"customerEmail" binding:"required,email"`
	Item          string `json:"item" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PENDING CONFIRMED SHIPPED DELIVERED CANCELLED"`
}
