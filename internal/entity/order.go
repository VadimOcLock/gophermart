package entity

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          uint64
	UserID      uint64
	OrderNumber string
	Status      OrderStatus
	Accrual     float64
	UploadedAt  time.Time
}
