package entity

import (
	"encoding/json"
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
	ID          uint64      `json:"-"`
	UserID      uint64      `json:"-"`
	OrderNumber string      `json:"number"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual,omitempty"`
	UploadedAt  time.Time   `json:"uploaded_at"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	alias := struct {
		*Alias
		UploadedAt string `json:"uploaded_at"`
	}{
		Alias:      (*Alias)(&o),
		UploadedAt: o.UploadedAt.Local().Format(time.RFC3339),
	}

	return json.Marshal(alias)
}
