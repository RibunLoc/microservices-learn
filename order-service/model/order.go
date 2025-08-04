package model

import (
	"time"

	"github.com/RibunLoc/microservices-learn/util"
	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64           `json:"order_id"`
	CustomerID  uuid.UUID        `json:"customer_id"`
	LineItems   []LineItem       `json:"Line_items"`
	OrderStatus string           `json:"order_status"`
	CreateAt    *time.Time       `json:"created_at"`
	ShippedAt   *util.CustomTime `json:"shipped_at,omitempty"`
	CompletedAt *util.CustomTime `json:"completed_at,omitempty"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
}
