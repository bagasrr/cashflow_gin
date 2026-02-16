package request

import (
	"time"
)

type CreateTransactionRequest struct {
	WalletID     string    `json:"wallet_id" binding:"required,uuid"`
	CategoryName string    `json:"category_name" binding:"required,max=100"`
	Title        string    `json:"title" binding:"required,max=255"`
	Amount       float64   `json:"amount" binding:"required,gt=0"` // Amount harus > 0
	Description  string    `json:"description"`
	Date         time.Time `json:"date" binding:"required"` // Format: RFC3339 (e.g., "2026-02-02T15:04:05Z")
}

// Untuk Update, biasanya field-nya optional (pake pointer)
type UpdateTransactionRequest struct {
	Title       string    `json:"title" binding:"omitempty,max=255"`
	Amount      float64   `json:"amount" binding:"omitempty,gt=0"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id" binding:"omitempty,uuid"`
	Date        time.Time `json:"date"`
}
