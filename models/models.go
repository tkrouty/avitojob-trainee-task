package models

import (
	"time"
)

type Transaction struct {
	TransactionID   int       `json:"transaction_id"`
	SourceID        string       `json:"source_id"`
	TargetID        string       `json:"target_id"`
	Sum             float64   `json:"sum" binding:"required"`
	TransactionTime time.Time `json:"transaction_time"`
}

type User struct {
	UserID  string   `json:"user_id"`
	Balance float64 `json:"balance"`
}
