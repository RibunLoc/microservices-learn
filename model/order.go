package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uuid.UUID  `json:"customer_id`
	LineItem    []LineItem `json:"Line_items`
	OrderStatus string     `json:"created-at`
	CreateAt    *time.Time `json:"shipped_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemID   uuid.UUID
	Quantity uint
	Price    uint
}
