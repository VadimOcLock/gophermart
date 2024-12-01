package entity

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	ID          uint64    `json:"-"`
	UserID      uint64    `json:"-"`
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

func (o Withdraw) MarshalJSON() ([]byte, error) {
	type Alias Withdraw
	alias := struct {
		*Alias
		ProcessedAt string `json:"processed_at"`
	}{
		Alias:       (*Alias)(&o),
		ProcessedAt: o.ProcessedAt.Local().Format(time.RFC3339),
	}

	return json.Marshal(alias)
}
